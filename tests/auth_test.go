package tests

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"mime/quotedprintable"
	"net"
	"net/http"
	"net/http/httptest"
	stdmail "net/mail"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"abclinic.local/backend/internal/auth"
	"abclinic.local/backend/internal/config"
	"abclinic.local/backend/internal/platform/clock"
	"abclinic.local/backend/internal/platform/mailer"
	"abclinic.local/backend/internal/platform/middleware"
	"abclinic.local/backend/migrations"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type memoryRepository struct {
	mu          sync.Mutex
	accounts    map[string]auth.Account
	revoked     map[string]bool
	created     auth.NewAccount
	emailTokens map[string]auth.EmailToken
	sender      *recordingSender
}

func (r *memoryRepository) Create(_ context.Context, value auth.NewAccount) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.accounts {
		if existing.Email == value.Email {
			return auth.Account{}, auth.ErrEmailExists
		}
	}
	r.accounts[value.ID] = value.Account
	r.created = value
	r.emailTokens[value.ID] = value.Verification
	return value.Account, nil
}
func (r *memoryRepository) FindByEmail(_ context.Context, email string) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, account := range r.accounts {
		if account.Email == email {
			return account, nil
		}
	}
	return auth.Account{}, auth.ErrNotFound
}
func (r *memoryRepository) FindByID(_ context.Context, id string) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	account, ok := r.accounts[id]
	if !ok {
		return auth.Account{}, auth.ErrNotFound
	}
	return account, nil
}
func (r *memoryRepository) IsInvalidated(_ context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.revoked[id], nil
}
func (r *memoryRepository) Revoke(_ context.Context, tokens []auth.InvalidatedToken, singleUse bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if singleUse && r.revoked[tokens[0].JTI] {
		return auth.ErrInvalidToken
	}
	for _, token := range tokens {
		r.revoked[token.JTI] = true
	}
	return nil
}

type testClock struct{ now time.Time }

func (c *testClock) Now() time.Time { return c.now }

func tokenConfig() auth.TokenConfig {
	return auth.TokenConfig{AccessSecret: strings.Repeat("a", 32), RefreshSecret: strings.Repeat("r", 32), Issuer: "test", AccessTTL: time.Minute, RefreshTTL: time.Hour}
}
func newTestService(t *testing.T) (*auth.Service, *memoryRepository, *testClock) {
	t.Helper()
	repo := &memoryRepository{accounts: map[string]auth.Account{}, revoked: map[string]bool{}, emailTokens: map[string]auth.EmailToken{}, sender: &recordingSender{}}
	c := &testClock{now: time.Now().Truncate(time.Second)}
	tokens, err := auth.NewTokens(tokenConfig(), c)
	if err != nil {
		t.Fatal(err)
	}
	verification, err := auth.NewEmailVerification(repo, repo.sender, c, emailConfig())
	if err != nil {
		t.Fatal(err)
	}
	service := auth.NewService(repo, tokens, verification)
	testSenders.Store(service, repo.sender)
	t.Cleanup(func() { testSenders.Delete(service) })
	return service, repo, c
}

var testSenders sync.Map

func registerUser(t *testing.T, s *auth.Service, email string) auth.UserResponse {
	t.Helper()
	user, err := s.Register(context.Background(), auth.RegisterRequest{Email: email, Password: " Password123! ", FullName: " Person "})
	if err != nil {
		t.Fatal(err)
	}
	verifyLastEmail(t, s, senderForService(t, s, email))
	return user
}
func loginUser(t *testing.T, s *auth.Service, email string) auth.LoginResponse {
	t.Helper()
	result, err := s.Login(context.Background(), auth.LoginRequest{Email: email, Password: " Password123! "})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func tokenRequest(pair auth.TokenResponse) auth.TokenRequest {
	return auth.TokenRequest{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}
}

func TestRegisterNormalizesEmailAndHashesPassword(t *testing.T) {
	s, repo, _ := newTestService(t)
	user := registerUser(t, s, " Person@Example.com ")
	if user.Email != "person@example.com" || user.FullName != "Person" || user.Role != auth.RolePatient {
		t.Fatalf("unexpected user: %+v", user)
	}
	if repo.created.PasswordHash == " Password123! " || bcrypt.CompareHashAndPassword([]byte(repo.created.PasswordHash), []byte(" Password123! ")) != nil {
		t.Fatal("password was not hashed correctly")
	}
	_, err := s.Register(context.Background(), auth.RegisterRequest{Email: "PERSON@example.com", Password: "Another123!", FullName: "Other"})
	if !errors.Is(err, auth.ErrEmailExists) {
		t.Fatalf("duplicate email: %v", err)
	}
	_, err = s.RegisterPatient(context.Background(), auth.RegisterPatientRequest{RegisterRequest: auth.RegisterRequest{Email: "patient@example.com", Password: "Password123!", FullName: "Patient"}, Phone: "0901234567"})
	if err != nil || repo.created.PatientID == "" || repo.created.Phone != "0901234567" {
		t.Fatalf("patient registration: %+v %v", repo.created, err)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	s, repo, _ := newTestService(t)
	for _, req := range []auth.RegisterRequest{
		{Email: "invalid", Password: "Password123!", FullName: "Person"},
		{Email: "Name <a@example.com>", Password: "Password123!", FullName: "Person"},
		{Email: "a@example.com", Password: "Password123!", FullName: " "},
	} {
		if _, err := s.Register(context.Background(), req); !errors.Is(err, auth.ErrInvalidInput) {
			t.Errorf("expected invalid input: %v", err)
		}
	}
	if _, err := s.RegisterPatient(context.Background(), auth.RegisterPatientRequest{Phone: "+84901234567"}); !errors.Is(err, auth.ErrInvalidInput) {
		t.Fatal("invalid phone accepted")
	}
	if len(repo.accounts) != 0 {
		t.Fatal("invalid registration must not write accounts")
	}
}

func TestLoginAndVerifyCurrentAccount(t *testing.T) {
	s, repo, _ := newTestService(t)
	user := registerUser(t, s, "person@example.com")
	for _, req := range []auth.LoginRequest{{Email: user.Email, Password: "wrong123"}, {Email: "missing@example.com", Password: "wrong123"}} {
		if _, err := s.Login(context.Background(), req); !errors.Is(err, auth.ErrInvalidLogin) {
			t.Fatalf("credentials: %v", err)
		}
	}
	login := loginUser(t, s, "PERSON@example.com")
	if login.RefreshToken == "" || login.AccessToken == "" || login.User.UserID != user.UserID {
		t.Fatal("missing token pair or user")
	}
	account := repo.accounts[user.UserID]
	account.Role = auth.RoleDoctor
	repo.accounts[user.UserID] = account
	principal, err := s.Verify(context.Background(), login.AccessToken)
	if err != nil || principal.Roles[0] != "DOCTOR" {
		t.Fatalf("current account role not used: %+v %v", principal, err)
	}
	account.Active = false
	repo.accounts[user.UserID] = account
	if _, err := s.Verify(context.Background(), login.AccessToken); !errors.Is(err, auth.ErrInactive) {
		t.Fatalf("inactive account: %v", err)
	}
	if _, err := s.Login(context.Background(), auth.LoginRequest{Email: user.Email, Password: " Password123! "}); !errors.Is(err, auth.ErrInactive) {
		t.Fatal("inactive account logged in")
	}
	delete(repo.accounts, user.UserID)
	if _, err := s.Verify(context.Background(), login.AccessToken); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("deleted account: %v", err)
	}
}

func TestRefreshRotationAndLogout(t *testing.T) {
	s, _, c := newTestService(t)
	registerUser(t, s, "person@example.com")
	login := loginUser(t, s, "person@example.com")
	c.now = c.now.Add(2 * time.Minute)
	if _, err := s.Verify(context.Background(), login.AccessToken); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("expired access token accepted")
	}
	refreshed, err := s.RefreshToken(context.Background(), tokenRequest(login.TokenResponse))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Verify(context.Background(), refreshed.AccessToken); err != nil {
		t.Fatal(err)
	}
	if _, err := s.RefreshToken(context.Background(), tokenRequest(login.TokenResponse)); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("reused refresh token: %v", err)
	}
	for range 2 {
		if err := s.Logout(context.Background(), tokenRequest(refreshed)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.Verify(context.Background(), refreshed.AccessToken); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("logged out token accepted")
	}
	if _, err := s.RefreshToken(context.Background(), tokenRequest(refreshed)); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("logged out refresh token accepted")
	}
}

func TestRefreshRejectsMixedUsersAndTokenTypes(t *testing.T) {
	s, _, _ := newTestService(t)
	registerUser(t, s, "first@example.com")
	registerUser(t, s, "second@example.com")
	first := loginUser(t, s, "first@example.com")
	second := loginUser(t, s, "second@example.com")
	for _, req := range []auth.TokenRequest{
		{AccessToken: first.AccessToken, RefreshToken: second.RefreshToken},
		{AccessToken: first.RefreshToken, RefreshToken: first.AccessToken},
	} {
		if _, err := s.RefreshToken(context.Background(), req); !errors.Is(err, auth.ErrInvalidToken) {
			t.Fatalf("invalid token pair: %v", err)
		}
		if err := s.Logout(context.Background(), req); !errors.Is(err, auth.ErrInvalidToken) {
			t.Fatalf("invalid logout token pair: %v", err)
		}
	}
	if _, err := s.RefreshToken(context.Background(), auth.TokenRequest{}); !errors.Is(err, auth.ErrInvalidInput) {
		t.Fatal("empty token pair accepted")
	}
}

func TestRefreshIsSingleUseUnderConcurrency(t *testing.T) {
	s, _, _ := newTestService(t)
	registerUser(t, s, "person@example.com")
	login := loginUser(t, s, "person@example.com")
	results := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := s.RefreshToken(context.Background(), tokenRequest(login.TokenResponse))
			results <- err
		}()
	}
	success, rejected := 0, 0
	for range 2 {
		switch err := <-results; {
		case err == nil:
			success++
		case errors.Is(err, auth.ErrInvalidToken):
			rejected++
		default:
			t.Fatal(err)
		}
	}
	if success != 1 || rejected != 1 {
		t.Fatalf("refresh outcomes: %d/%d", success, rejected)
	}
}

func TestTokensRejectInvalidClaimsAndSignatures(t *testing.T) {
	s, _, c := newTestService(t)
	user := registerUser(t, s, "person@example.com")
	now := c.Now()
	original := jwt.MapClaims{"jti": "id", "sub": user.UserID, "iss": "test", "iat": now.Unix(), "exp": now.Add(time.Minute).Unix()}
	cases := []struct {
		name   string
		change func(jwt.MapClaims)
		method jwt.SigningMethod
		secret string
	}{
		{name: "expired", change: func(c jwt.MapClaims) { c["exp"] = now.Add(-time.Second).Unix() }},
		{name: "missing expiry", change: func(c jwt.MapClaims) { delete(c, "exp") }},
		{name: "missing subject", change: func(c jwt.MapClaims) { delete(c, "sub") }},
		{name: "missing ID", change: func(c jwt.MapClaims) { delete(c, "jti") }},
		{name: "wrong issuer", change: func(c jwt.MapClaims) { c["iss"] = "other" }},
		{name: "future issued at", change: func(c jwt.MapClaims) { c["iat"] = now.Add(time.Hour).Unix() }},
		{name: "missing issued at", change: func(c jwt.MapClaims) { delete(c, "iat") }},
		{name: "future not before", change: func(c jwt.MapClaims) { c["nbf"] = now.Add(time.Hour).Unix() }},
		{name: "wrong signature", secret: strings.Repeat("x", 32)},
		{name: "wrong algorithm", method: jwt.SigningMethodHS384},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value := jwt.MapClaims{}
			for k, v := range original {
				value[k] = v
			}
			if tc.change != nil {
				tc.change(value)
			}
			method, secret := tc.method, tc.secret
			if method == nil {
				method = jwt.SigningMethodHS256
			}
			if secret == "" {
				secret = tokenConfig().AccessSecret
			}
			raw, err := jwt.NewWithClaims(method, value).SignedString([]byte(secret))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := s.Verify(context.Background(), raw); !errors.Is(err, auth.ErrInvalidToken) {
				t.Fatalf("expected invalid token: %v", err)
			}
		})
	}
	if _, err := s.Verify(context.Background(), "malformed"); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("malformed token accepted")
	}
}

func authRouter(t *testing.T) *gin.Engine {
	router, _ := authRouterWithState(t)
	return router
}
func authRouterWithState(t *testing.T) (*gin.Engine, *memoryRepository) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	s, repo, _ := newTestService(t)
	router := gin.New()
	auth.NewHandler(s).RegisterRoutes(router.Group("/api"))
	router.GET("/protected", middleware.Authenticate(s), func(c *gin.Context) { c.Status(200) })
	return router, repo
}
func request(router http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
func TestHTTPAuthenticationFlow(t *testing.T) {
	router, repo := authRouterWithState(t)
	body := `{"email":"person@example.com","password":"Password123!","fullName":"Person","phone":"0901234567"}`
	registered := request(router, "POST", "/api/auth/register-patient", body, "")
	if registered.Code != 201 {
		t.Fatalf("registration: %d %s", registered.Code, registered.Body)
	}
	if strings.Contains(registered.Body.String(), "password") {
		t.Fatal("password leaked")
	}
	blocked := request(router, "POST", "/api/auth/login", `{"email":"person@example.com","password":"Password123!"}`, "")
	if blocked.Code != 403 {
		t.Fatalf("unverified account logged in: %d", blocked.Code)
	}
	verificationBody, _ := json.Marshal(auth.VerifyEmailRequest{Token: repo.sender.token(t, "person@example.com")})
	verified := request(router, "POST", "/api/auth/verify-email", string(verificationBody), "")
	if verified.Code != 200 {
		t.Fatalf("verify email: %d %s", verified.Code, verified.Body)
	}
	login := request(router, "POST", "/api/auth/login", `{"email":"person@example.com","password":"Password123!"}`, "")
	var response struct {
		StatusCode int                `json:"statusCode"`
		Data       auth.LoginResponse `json:"data"`
	}
	if err := json.Unmarshal(login.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if login.Code != 200 || response.StatusCode != 1000 || response.Data.RefreshToken == "" {
		t.Fatalf("login: %d %s", login.Code, login.Body)
	}
	if protected := request(router, "GET", "/protected", "", response.Data.AccessToken); protected.Code != 200 {
		t.Fatal("valid JWT rejected by existing middleware")
	}
	raw, _ := json.Marshal(tokenRequest(response.Data.TokenResponse))
	refresh := request(router, "POST", "/api/auth/refresh_token", string(raw), "")
	if refresh.Code != 200 {
		t.Fatalf("refresh: %d %s", refresh.Code, refresh.Body)
	}
	var pair struct {
		Data auth.TokenResponse `json:"data"`
	}
	if err := json.Unmarshal(refresh.Body.Bytes(), &pair); err != nil {
		t.Fatal(err)
	}
	logoutBody, _ := json.Marshal(tokenRequest(pair.Data))
	logout := request(router, "POST", "/api/auth/logout", string(logoutBody), "")
	if logout.Code != 200 {
		t.Fatalf("logout: %d %s", logout.Code, logout.Body)
	}
	if protected := request(router, "GET", "/protected", "", pair.Data.AccessToken); protected.Code != 401 {
		t.Fatal("revoked JWT accepted")
	}
	if login.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("token response is cacheable")
	}
}
func TestHTTPRejectsMalformedAndPrivilegedRegistration(t *testing.T) {
	router := authRouter(t)
	for _, body := range []string{
		`{`, `null`,
		`{"email":"a@example.com","password":"Password123!","fullName":"Person","role":"ADMIN"}`,
		`{"email":"a@example.com","password":"Password123!","fullName":"Person"} {}`,
		`{"email":"` + strings.Repeat("x", 17<<10) + `"}`,
	} {
		result := request(router, "POST", "/api/auth/register", body, "")
		if result.Code != 400 {
			t.Fatalf("invalid request: %d %s", result.Code, result.Body)
		}
	}
}

func TestAuthPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL auth integration tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	control, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal("cannot connect to test PostgreSQL")
	}
	t.Cleanup(func() { _ = control.Close(context.Background()) })
	var random [12]byte
	if _, err := rand.Read(random[:]); err != nil {
		t.Fatal(err)
	}
	schemaName := "auth_test_" + hex.EncodeToString(random[:])
	schema := pgx.Identifier{schemaName}.Sanitize()
	if _, err := control.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := control.Exec(cleanup, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Error(err)
		}
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid test PostgreSQL URL")
	}
	cfg.MaxConns = 4
	cfg.ConnConfig.RuntimeParams["search_path"] = schemaName
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	for range 2 {
		if err := migrations.Apply(ctx, pool); err != nil {
			t.Fatalf("apply migration: %v", err)
		}
	}
	tokens, err := auth.NewTokens(tokenConfig(), clock.Fixed{Time: time.Now().Truncate(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	repo := auth.NewPostgresRepository(pool)
	sender := &recordingSender{}
	verification, err := auth.NewEmailVerification(repo, sender, clock.Fixed{Time: time.Now().Truncate(time.Second)}, emailConfig())
	if err != nil {
		t.Fatal(err)
	}
	service := auth.NewService(repo, tokens, verification)
	if _, err := pool.Exec(ctx, `INSERT INTO patients(id,full_name,phone) VALUES('existing','Existing Patient','0901111111')`); err != nil {
		t.Fatal(err)
	}
	user, err := service.RegisterPatient(ctx, auth.RegisterPatientRequest{RegisterRequest: auth.RegisterRequest{Email: "person@example.com", Password: "Password123!", FullName: "Person"}, Phone: "0901111111"})
	if err != nil {
		t.Fatal(err)
	}
	var linked string
	if err := pool.QueryRow(ctx, `SELECT user_id FROM patients WHERE id='existing'`).Scan(&linked); err != nil || linked != user.UserID {
		t.Fatalf("existing patient not linked: %v", err)
	}
	_, err = service.RegisterPatient(ctx, auth.RegisterPatientRequest{RegisterRequest: auth.RegisterRequest{Email: "other@example.com", Password: "Password123!", FullName: "Other"}, Phone: "0901111111"})
	if !errors.Is(err, auth.ErrPhoneExists) {
		t.Fatalf("duplicate phone: %v", err)
	}
	if _, err := repo.FindByEmail(ctx, "other@example.com"); !errors.Is(err, auth.ErrNotFound) {
		t.Fatal("failed registration left a user behind")
	}
	if _, err := pool.Exec(ctx, `INSERT INTO patients(id,phone) VALUES('ambiguous1','0902222222'),('ambiguous2','0902222222')`); err != nil {
		t.Fatal(err)
	}
	_, err = service.RegisterPatient(ctx, auth.RegisterPatientRequest{RegisterRequest: auth.RegisterRequest{Email: "ambiguous@example.com", Password: "Password123!", FullName: "Person"}, Phone: "0902222222"})
	if !errors.Is(err, auth.ErrPatientConflict) {
		t.Fatalf("ambiguous patient: %v", err)
	}
	if _, err := repo.FindByEmail(ctx, "ambiguous@example.com"); !errors.Is(err, auth.ErrNotFound) {
		t.Fatal("ambiguous profile registration did not roll back")
	}
	results := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := service.Register(ctx, auth.RegisterRequest{Email: "concurrent@example.com", Password: "Password123!", FullName: "Person"})
			results <- err
		}()
	}
	success, duplicates := 0, 0
	for range 2 {
		switch err := <-results; {
		case err == nil:
			success++
		case errors.Is(err, auth.ErrEmailExists):
			duplicates++
		default:
			t.Fatal(err)
		}
	}
	if success != 1 || duplicates != 1 {
		t.Fatalf("concurrent registration: %d/%d", success, duplicates)
	}
	verifyLastEmail(t, service, sender.token(t, "person@example.com"))
	login, err := service.Login(ctx, auth.LoginRequest{Email: "PERSON@example.com", Password: "Password123!"})
	if err != nil {
		t.Fatal(err)
	}
	pair, err := service.RefreshToken(ctx, tokenRequest(login.TokenResponse))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.RefreshToken(ctx, tokenRequest(login.TokenResponse)); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("refresh token reused")
	}
	if err := service.Logout(ctx, tokenRequest(pair)); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Verify(ctx, pair.AccessToken); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("logged out JWT accepted")
	}

	extra, err := service.Register(ctx, auth.RegisterRequest{Email: "sql-verify@example.com", Password: "Password123!", FullName: "SQL Verify"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Login(ctx, auth.LoginRequest{Email: extra.Email, Password: "Password123!"}); !errors.Is(err, auth.ErrEmailUnverified) {
		t.Fatal("SQL account logged in before email verification")
	}
	oldToken := sender.token(t, extra.Email)
	if _, err := pool.Exec(ctx, `UPDATE email_verifications SET expires_at=$2,last_sent_at=$3 WHERE user_id=$1`, extra.UserID, time.Now().Add(-time.Minute), time.Now().Add(-2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := service.VerifyEmail(ctx, auth.VerifyEmailRequest{Token: oldToken}); !errors.Is(err, auth.ErrInvalidEmailToken) {
		t.Fatal("SQL accepted expired email token")
	}
	if err := service.ResendVerification(ctx, auth.ResendVerificationRequest{Email: extra.Email}); err != nil {
		t.Fatal(err)
	}
	newToken := sender.token(t, extra.Email)
	if newToken == oldToken {
		t.Fatal("SQL resend did not replace email token")
	}
	if err := service.VerifyEmail(ctx, auth.VerifyEmailRequest{Token: oldToken}); !errors.Is(err, auth.ErrInvalidEmailToken) {
		t.Fatal("SQL accepted replaced email token")
	}
	verifyLastEmail(t, service, newToken)
	if err := service.VerifyEmail(ctx, auth.VerifyEmailRequest{Token: newToken}); !errors.Is(err, auth.ErrInvalidEmailToken) {
		t.Fatal("SQL reused email token")
	}
}

func TestSwaggerSameOriginAuth(t *testing.T) {
	for _, tc := range []struct {
		name, target, origin string
		status               int
	}{
		{"swagger localhost", "http://localhost:8081/api/auth/register", "http://localhost:8081", 201},
		{"swagger loopback", "http://127.0.0.1:8081/api/auth/register", "http://127.0.0.1:8081", 201},
		{"swagger HTTPS", "https://api.example.com/api/auth/register", "https://api.example.com", 201},
		{"configured frontend", "http://localhost:8081/api/auth/register", "http://localhost:5173", 201},
		{"unconfigured website", "http://localhost:8081/api/auth/register", "https://other.example", 403},
		{"different scheme", "http://localhost:8081/api/auth/register", "https://localhost:8081", 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			service, _, _ := newTestService(t)
			router := gin.New()
			router.Use(middleware.CORS([]string{"http://localhost:5173", "http://localhost:3000"}))
			auth.NewHandler(service).RegisterRoutes(router.Group("/api"))
			req := httptest.NewRequest(http.MethodPost, tc.target, strings.NewReader(`{"email":"swagger@example.com","password":"Password123!","fullName":"Swagger Test"}`))
			req.Header.Set("Origin", tc.origin)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("got %d, want %d: %s", rec.Code, tc.status, rec.Body)
			}
			if tc.status == 201 && rec.Header().Get("Access-Control-Allow-Origin") != tc.origin {
				t.Fatal("missing allowed origin")
			}
		})
	}
}

func TestRegistrationPasswordPolicy(t *testing.T) {
	cases := []struct {
		name, password string
		valid          bool
	}{
		{"minimum length", "Aa1!xxxx", true},
		{"maximum length", "Aa1!" + strings.Repeat("x", 28), true},
		{"Unicode letters", "Áa1!" + strings.Repeat("界", 4), true},
		{"too short", "Aa1!xxx", false},
		{"too long", "Aa1!" + strings.Repeat("x", 29), false},
		{"missing uppercase", "password123!", false},
		{"missing lowercase", "PASSWORD123!", false},
		{"missing number", "Password!!!", false},
		{"missing special character", "Password123", false},
		{"space is not special", "Password123 ", false},
		{"empty password", "", false},
		{"invalid UTF-8", "Aa1!xxxx" + string([]byte{0xff}), false},
		{"bcrypt byte limit", "Aa1!" + strings.Repeat("界", 28), false},
	}
	for _, tc := range cases {
		for _, patient := range []bool{false, true} {
			endpoint := "register"
			if patient {
				endpoint = "register-patient"
			}
			t.Run(endpoint+"/"+tc.name, func(t *testing.T) {
				service, repo, _ := newTestService(t)
				request := auth.RegisterRequest{Email: "policy@example.com", Password: tc.password, FullName: "Policy Test"}
				var err error
				if patient {
					_, err = service.RegisterPatient(context.Background(), auth.RegisterPatientRequest{RegisterRequest: request, Phone: "0901234567"})
				} else {
					_, err = service.Register(context.Background(), request)
				}
				if tc.valid {
					if err != nil {
						t.Fatal(err)
					}
					if len(repo.accounts) != 1 {
						t.Fatal("valid password did not create an account")
					}
				} else {
					if !errors.Is(err, auth.ErrWeakPassword) {
						t.Fatalf("expected password policy error, got %v", err)
					}
					if len(repo.accounts) != 0 {
						t.Fatal("weak password created an account")
					}
				}
			})
		}
	}
}

func TestHTTPWeakPasswordReturnsPolicyMessage(t *testing.T) {
	router := authRouter(t)
	result := request(router, "POST", "/api/auth/register", `{"email":"policy@example.com","password":"password123","fullName":"Person"}`, "")
	var body struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}
	if err := json.Unmarshal(result.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if result.Code != http.StatusBadRequest || body.StatusCode != 400 || body.Message != auth.ErrWeakPassword.Message {
		t.Fatalf("unexpected password policy response: %d %s", result.Code, result.Body)
	}
}

func TestExistingPasswordStillLogsIn(t *testing.T) {
	service, repo, _ := newTestService(t)
	user := registerUser(t, service, "legacy@example.com")
	legacyPassword := "legacy123"
	hash, err := bcrypt.GenerateFromPassword([]byte(legacyPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	account := repo.accounts[user.UserID]
	account.PasswordHash = string(hash)
	repo.accounts[user.UserID] = account
	if _, err := service.Login(context.Background(), auth.LoginRequest{Email: user.Email, Password: legacyPassword}); err != nil {
		t.Fatalf("password creation policy must not block an existing login: %v", err)
	}
}

type sentEmail struct{ to, subject, body string }

type recordingSender struct {
	mu       sync.Mutex
	messages []sentEmail
	failure  error
}

func (s *recordingSender) Send(_ context.Context, to, subject, body string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failure != nil {
		return s.failure
	}
	s.messages = append(s.messages, sentEmail{to: to, subject: subject, body: body})
	return nil
}
func (s *recordingSender) token(t *testing.T, email string) string {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := len(s.messages) - 1; i >= 0; i-- {
		message := s.messages[i]
		if message.to != normalizeTestEmail(email) {
			continue
		}
		for _, word := range strings.Fields(message.body) {
			link, err := url.Parse(word)
			if err != nil || link.Fragment == "" {
				continue
			}
			values, err := url.ParseQuery(link.Fragment)
			if err == nil && values.Get("token") != "" {
				return values.Get("token")
			}
		}
	}
	t.Fatal("verification email was not captured")
	return ""
}
func (s *recordingSender) count() int        { s.mu.Lock(); defer s.mu.Unlock(); return len(s.messages) }
func normalizeTestEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }
func senderForService(t *testing.T, service *auth.Service, email string) string {
	t.Helper()
	value, ok := testSenders.Load(service)
	if !ok {
		t.Fatal("missing test sender")
	}
	return value.(*recordingSender).token(t, email)
}
func emailConfig() auth.EmailConfig {
	return auth.EmailConfig{VerificationURL: "http://localhost:8081/api/auth/verify-email", TokenTTL: 15 * time.Minute, ResendCooldown: time.Minute}
}
func verifyLastEmail(t *testing.T, service *auth.Service, token string) {
	t.Helper()
	if err := service.VerifyEmail(context.Background(), auth.VerifyEmailRequest{Token: token}); err != nil {
		t.Fatal(err)
	}
}

func (r *memoryRepository) ReplaceEmailToken(_ context.Context, token auth.EmailToken, cutoff time.Time) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	account, ok := r.accounts[token.AccountID]
	if !ok || !account.Active || account.EmailVerified {
		return false, nil
	}
	if previous, ok := r.emailTokens[token.AccountID]; ok && previous.SentAt.After(cutoff) {
		return false, nil
	}
	r.emailTokens[token.AccountID] = token
	return true, nil
}
func (r *memoryRepository) VerifyEmail(_ context.Context, hash string, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, token := range r.emailTokens {
		if token.Hash != hash {
			continue
		}
		account, ok := r.accounts[id]
		if !ok || account.EmailVerified || !token.ExpiresAt.After(now) {
			return auth.ErrInvalidEmailToken
		}
		if !account.Active {
			return auth.ErrInactive
		}
		account.EmailVerified = true
		r.accounts[id] = account
		delete(r.emailTokens, id)
		return nil
	}
	return auth.ErrInvalidEmailToken
}

func TestEmailVerificationRequiredAndSingleUse(t *testing.T) {
	service, repo, _ := newTestService(t)
	user, err := service.Register(context.Background(), auth.RegisterRequest{Email: "verify@example.com", Password: "Password123!", FullName: "Person"})
	if err != nil {
		t.Fatal(err)
	}
	if user.EmailVerified {
		t.Fatal("registration marked email verified")
	}
	if _, err := service.Login(context.Background(), auth.LoginRequest{Email: user.Email, Password: "Password123!"}); !errors.Is(err, auth.ErrEmailUnverified) {
		t.Fatal("unverified email was able to log in")
	}
	token := repo.sender.token(t, user.Email)
	decoded, _ := base64.RawURLEncoding.DecodeString(token)
	digest := sha256.Sum256(decoded)
	stored := repo.emailTokens[user.UserID]
	if stored.Hash != hex.EncodeToString(digest[:]) || strings.Contains(stored.Hash, token) {
		t.Fatal("verification token was not hashed")
	}
	verifyLastEmail(t, service, token)
	if _, err := service.Login(context.Background(), auth.LoginRequest{Email: user.Email, Password: "Password123!"}); err != nil {
		t.Fatal(err)
	}
	if err := service.VerifyEmail(context.Background(), auth.VerifyEmailRequest{Token: token}); !errors.Is(err, auth.ErrInvalidEmailToken) {
		t.Fatal("verification token was reused")
	}
}

func TestExpiredEmailAndResend(t *testing.T) {
	service, repo, c := newTestService(t)
	user, err := service.Register(context.Background(), auth.RegisterRequest{Email: "resend@example.com", Password: "Password123!", FullName: "Person"})
	if err != nil {
		t.Fatal(err)
	}
	initial := repo.sender.token(t, user.Email)
	for _, email := range []string{user.Email, "missing@example.com"} {
		if err := service.ResendVerification(context.Background(), auth.ResendVerificationRequest{Email: email}); err != nil {
			t.Fatal(err)
		}
	}
	if repo.sender.count() != 1 {
		t.Fatal("resend ignored cooldown or sent to unknown account")
	}
	c.now = c.now.Add(15 * time.Minute)
	if err := service.VerifyEmail(context.Background(), auth.VerifyEmailRequest{Token: initial}); !errors.Is(err, auth.ErrInvalidEmailToken) {
		t.Fatal("expired verification token accepted")
	}
	if err := service.ResendVerification(context.Background(), auth.ResendVerificationRequest{Email: user.Email}); err != nil {
		t.Fatal(err)
	}
	replacement := repo.sender.token(t, user.Email)
	if replacement == initial {
		t.Fatal("resend reused old token")
	}
	if err := service.VerifyEmail(context.Background(), auth.VerifyEmailRequest{Token: initial}); !errors.Is(err, auth.ErrInvalidEmailToken) {
		t.Fatal("old token survived resend")
	}
	verifyLastEmail(t, service, replacement)
	if err := service.ResendVerification(context.Background(), auth.ResendVerificationRequest{Email: user.Email}); err != nil {
		t.Fatal(err)
	}
	if repo.sender.count() != 2 {
		t.Fatal("sent verification to already verified account")
	}
}

func TestSMTPFailureLeavesRecoverablePendingAccount(t *testing.T) {
	service, repo, c := newTestService(t)
	repo.sender.failure = errors.New("SMTP unavailable")
	_, err := service.Register(context.Background(), auth.RegisterRequest{Email: "retry@example.com", Password: "Password123!", FullName: "Person"})
	if !errors.Is(err, auth.ErrEmailDelivery) || len(repo.accounts) != 1 {
		t.Fatalf("expected pending account and delivery error, got %v", err)
	}
	repo.sender.failure = nil
	c.now = c.now.Add(time.Minute)
	if err := service.ResendVerification(context.Background(), auth.ResendVerificationRequest{Email: "retry@example.com"}); err != nil {
		t.Fatal(err)
	}
	verifyLastEmail(t, service, repo.sender.token(t, "retry@example.com"))
}

func TestRegistrationRequiresEmailSender(t *testing.T) {
	_, repo, _ := newTestService(t)
	tokens, err := auth.NewTokens(tokenConfig(), clock.Fixed{Time: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	unavailable := auth.NewService(repo, tokens, nil)
	_, err = unavailable.Register(context.Background(), auth.RegisterRequest{Email: "no-smtp@example.com", Password: "Password123!", FullName: "Person"})
	if !errors.Is(err, auth.ErrEmailNotConfigured) || len(repo.accounts) != 0 {
		t.Fatal("missing sender created an account")
	}
}

func TestEmailVerificationCannotUnblockInactiveAccount(t *testing.T) {
	service, repo, _ := newTestService(t)
	user, err := service.Register(context.Background(), auth.RegisterRequest{Email: "blocked@example.com", Password: "Password123!", FullName: "Person"})
	if err != nil {
		t.Fatal(err)
	}
	account := repo.accounts[user.UserID]
	account.Active = false
	repo.accounts[user.UserID] = account
	err = service.VerifyEmail(context.Background(), auth.VerifyEmailRequest{Token: repo.sender.token(t, user.Email)})
	if !errors.Is(err, auth.ErrInactive) || repo.accounts[user.UserID].Active {
		t.Fatal("email verification bypassed account lock")
	}
}

func TestVerificationPageDoesNotConsumeToken(t *testing.T) {
	router, repo := authRouterWithState(t)
	registered := request(router, "POST", "/api/auth/register", `{"email":"page@example.com","password":"Password123!","fullName":"Person"}`, "")
	if registered.Code != 201 {
		t.Fatal(registered.Body)
	}
	page := request(router, "GET", "/api/auth/verify-email", "", "")
	if page.Code != 200 || !strings.Contains(page.Header().Get("Content-Type"), "text/html") || page.Header().Get("Referrer-Policy") != "no-referrer" || !strings.Contains(page.Header().Get("Content-Security-Policy"), "nonce-") {
		t.Fatal("verification page is not protected")
	}
	if strings.Contains(page.Body.String(), repo.sender.token(t, "page@example.com")) {
		t.Fatal("verification token leaked into HTML")
	}
	login := request(router, "POST", "/api/auth/login", `{"email":"page@example.com","password":"Password123!"}`, "")
	if login.Code != 403 {
		t.Fatal("GET verification page activated account")
	}
	body, _ := json.Marshal(auth.VerifyEmailRequest{Token: repo.sender.token(t, "page@example.com")})
	if verified := request(router, "POST", "/api/auth/verify-email", string(body), ""); verified.Code != 200 {
		t.Fatal(verified.Body)
	}
}

func TestConcurrentEmailResendAndVerification(t *testing.T) {
	service, repo, c := newTestService(t)
	user, err := service.Register(context.Background(), auth.RegisterRequest{Email: "concurrent-email@example.com", Password: "Password123!", FullName: "Person"})
	if err != nil {
		t.Fatal(err)
	}
	c.now = c.now.Add(time.Minute)
	results := make(chan error, 2)
	for range 2 {
		go func() {
			results <- service.ResendVerification(context.Background(), auth.ResendVerificationRequest{Email: user.Email})
		}()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	if repo.sender.count() != 2 {
		t.Fatal("concurrent resend sent duplicate email")
	}
	token := repo.sender.token(t, user.Email)
	for range 2 {
		go func() { results <- service.VerifyEmail(context.Background(), auth.VerifyEmailRequest{Token: token}) }()
	}
	success, rejected := 0, 0
	for range 2 {
		switch err := <-results; {
		case err == nil:
			success++
		case errors.Is(err, auth.ErrInvalidEmailToken):
			rejected++
		default:
			t.Fatal(err)
		}
	}
	if success != 1 || rejected != 1 {
		t.Fatal("verification was not single-use under concurrency")
	}
}

func TestEmailSMTPConfiguration(t *testing.T) {
	keys := []string{"SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM", "SMTP_ENCRYPTION", "SMTP_TIMEOUT", "EMAIL_VERIFICATION_URL", "EMAIL_VERIFICATION_TTL", "EMAIL_RESEND_COOLDOWN"}
	for _, key := range keys {
		t.Setenv(key, "")
	}
	t.Setenv("SMTP_HOST", "smtp.gmail.com")
	t.Setenv("SMTP_USERNAME", "sender@example.com")
	t.Setenv("SMTP_PASSWORD", "test-app-password")
	t.Setenv("SMTP_FROM", "Booking App <sender@example.com>")
	if _, err := config.FromEnv(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ key, value string }{
		{"SMTP_PORT", "0"}, {"SMTP_PORT", "65536"}, {"SMTP_ENCRYPTION", "none"},
		{"SMTP_FROM", "invalid"}, {"SMTP_FROM", "sender@example.com\r\nBcc:other@example.com"},
		{"SMTP_TIMEOUT", "0s"}, {"EMAIL_VERIFICATION_URL", "javascript:alert(1)"},
		{"EMAIL_VERIFICATION_URL", "https://example.com/verify?token=unsafe"},
		{"EMAIL_VERIFICATION_TTL", "0s"}, {"EMAIL_RESEND_COOLDOWN", "30m"},
	} {
		t.Run(tc.key+tc.value, func(t *testing.T) {
			t.Setenv(tc.key, tc.value)
			if _, err := config.FromEnv(); err == nil {
				t.Fatal("invalid mail configuration accepted")
			}
		})
	}
}

type smtpDelivery struct {
	data          []byte
	authenticated bool
	err           error
}

func smtpTestServer(t *testing.T, encryption string, offerSTARTTLS bool) (mailer.Config, <-chan smtpDelivery) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	certificate := &x509.Certificate{SerialNumber: big.NewInt(now.UnixNano()), NotBefore: now.Add(-time.Minute), NotAfter: now.Add(time.Hour), IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)}, KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IsCA: true, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(parsed)
	serverCertificate := tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	delivery := make(chan smtpDelivery, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			delivery <- smtpDelivery{err: err}
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(5 * time.Second))
		result := serveSMTP(connection, serverCertificate, encryption, offerSTARTTLS)
		delivery <- result
	}()
	return mailer.Config{Host: "127.0.0.1", Port: listener.Addr().(*net.TCPAddr).Port, Username: "sender@example.com", Password: "test-password", From: "Booking App <sender@example.com>", Encryption: encryption, Timeout: 3 * time.Second, RootCAs: roots}, delivery
}

func serveSMTP(connection net.Conn, certificate tls.Certificate, encryption string, offerSTARTTLS bool) smtpDelivery {
	result := smtpDelivery{}
	encrypted := encryption == "tls"
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{certificate}, MinVersion: tls.VersionTLS12}
	if encrypted {
		secure := tls.Server(connection, tlsConfig)
		if err := secure.Handshake(); err != nil {
			result.err = err
			return result
		}
		connection = secure
	}
	protocol := textproto.NewConn(connection)
	if err := protocol.PrintfLine("220 local test SMTP"); err != nil {
		result.err = err
		return result
	}
	for {
		line, err := protocol.ReadLine()
		if err != nil {
			result.err = err
			return result
		}
		switch {
		case strings.HasPrefix(line, "EHLO "):
			if !encrypted && offerSTARTTLS {
				err = protocol.PrintfLine("250-localhost\r\n250 STARTTLS")
			} else {
				err = protocol.PrintfLine("250-localhost\r\n250 AUTH PLAIN")
			}
		case line == "STARTTLS":
			if !offerSTARTTLS {
				result.err = errors.New("unexpected STARTTLS")
				return result
			}
			if err = protocol.PrintfLine("220 begin TLS"); err != nil {
				result.err = err
				return result
			}
			secure := tls.Server(connection, tlsConfig)
			if err = secure.Handshake(); err != nil {
				result.err = err
				return result
			}
			connection = secure
			protocol = textproto.NewConn(connection)
			encrypted = true
		case strings.HasPrefix(line, "AUTH PLAIN "):
			credentials, decodeErr := base64.StdEncoding.DecodeString(strings.TrimPrefix(line, "AUTH PLAIN "))
			if !encrypted || decodeErr != nil || string(credentials) != "\x00sender@example.com\x00test-password" {
				result.err = errors.New("invalid or insecure SMTP authentication")
				return result
			}
			result.authenticated = true
			err = protocol.PrintfLine("235 authenticated")
		case strings.HasPrefix(line, "MAIL FROM:"), strings.HasPrefix(line, "RCPT TO:"):
			err = protocol.PrintfLine("250 OK")
		case line == "DATA":
			if err = protocol.PrintfLine("354 send message"); err != nil {
				result.err = err
				return result
			}
			result.data, err = protocol.ReadDotBytes()
			if err == nil {
				err = protocol.PrintfLine("250 accepted")
			}
		case line == "QUIT":
			result.err = protocol.PrintfLine("221 bye")
			return result
		default:
			result.err = fmt.Errorf("unexpected SMTP command")
			return result
		}
		if err != nil {
			result.err = err
			return result
		}
	}
}

func TestSMTPVerificationDeliveryUsesTLS(t *testing.T) {
	for _, encryption := range []string{"tls", "starttls"} {
		t.Run(encryption, func(t *testing.T) {
			cfg, delivered := smtpTestServer(t, encryption, true)
			sender, err := mailer.NewSMTP(cfg)
			if err != nil {
				t.Fatal(err)
			}
			body := "Xác thực email Booking App\nhttp://localhost:8081/api/auth/verify-email#token=local-test-only\n"
			if err := sender.Send(context.Background(), "recipient@example.com", "Xác thực email", body); err != nil {
				t.Fatal(err)
			}
			result := <-delivered
			if result.err != nil || !result.authenticated {
				t.Fatalf("SMTP delivery: %v", result.err)
			}
			message, err := stdmail.ReadMessage(strings.NewReader(string(result.data)))
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := io.ReadAll(quotedprintable.NewReader(message.Body))
			if err != nil {
				t.Fatal(err)
			}
			if strings.ReplaceAll(string(decoded), "\r\n", "\n") != body {
				t.Fatalf("SMTP message body: got %q, want %q", string(decoded), body)
			}
		})
	}
}

func TestSMTPRefusesPlaintextFallback(t *testing.T) {
	cfg, delivered := smtpTestServer(t, "starttls", false)
	sender, err := mailer.NewSMTP(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), "recipient@example.com", "Subject", "Body"); err == nil {
		t.Fatal("SMTP accepted a server without STARTTLS")
	}
	result := <-delivered
	if result.authenticated || len(result.data) != 0 {
		t.Fatal("credentials or email sent without TLS")
	}
}
