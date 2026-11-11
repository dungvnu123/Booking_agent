package tests

import (
	"abclinic.local/backend/internal/auth"
	"abclinic.local/backend/internal/platform/clock"
	"abclinic.local/backend/internal/platform/middleware"
	"abclinic.local/backend/migrations"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type memoryRepository struct {
	mu       sync.Mutex
	accounts map[string]auth.Account
	revoked  map[string]bool
	created  auth.NewAccount
}
type testClock struct{ now time.Time }

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

func (c *testClock) Now() time.Time { return c.now }

func tokenConfig() auth.TokenConfig {
	return auth.TokenConfig{AccessSecret: strings.Repeat("a", 32), RefreshSecret: strings.Repeat("r", 32), Issuer: "test", AccessTTL: time.Minute, RefreshTTL: time.Hour}
}

func newTestService(t *testing.T) (*auth.Service, *memoryRepository, *testClock) {
	t.Helper()
	repo := &memoryRepository{accounts: map[string]auth.Account{}, revoked: map[string]bool{}}
	c := &testClock{now: time.Now().Truncate(time.Second)}
	tokens, err := auth.NewTokens(tokenConfig(), c)
	if err != nil {
		t.Fatal(err)
	}
	return auth.NewService(repo, tokens), repo, c
}

func registerUser(t *testing.T, s *auth.Service, email string) auth.UserResponse {
	t.Helper()
	user, err := s.Register(context.Background(), auth.RegisterRequest{Email: email, Password: " Password123! ", FullName: " Person "})
	if err != nil {
		t.Fatal(err)
	}
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
	router := authRouter(t)
	body := `{"email":"person@example.com","password":"Password123!","fullName":"Person","phone":"0901234567"}`
	registered := request(router, "POST", "/api/auth/register-patient", body, "")
	if registered.Code != 201 {
		t.Fatalf("registration: %d %s", registered.Code, registered.Body)
	}
	if strings.Contains(registered.Body.String(), "password") {
		t.Fatal("password leaked")
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
	service := auth.NewService(repo, tokens)
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

}
