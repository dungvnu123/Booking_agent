package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"abclinic.local/backend/internal/platform/middleware"
	"golang.org/x/crypto/bcrypt"
)

var phonePattern = regexp.MustCompile(`^0[0-9]{9}$`)

type Service struct {
	repository Repository
	tokens     *Tokens
	email      *EmailVerification
}

func NewService(repository Repository, tokens *Tokens, email *EmailVerification) *Service {
	return &Service{repository: repository, tokens: tokens, email: email}
}

func (s *Service) Register(ctx context.Context, request RegisterRequest) (UserResponse, error) {
	return s.register(ctx, request, "")
}

func (s *Service) RegisterPatient(ctx context.Context, request RegisterPatientRequest) (UserResponse, error) {
	phone := strings.TrimSpace(request.Phone)
	if !phonePattern.MatchString(phone) {
		return UserResponse{}, ErrInvalidInput
	}
	return s.register(ctx, request.RegisterRequest, phone)
}

func (s *Service) register(ctx context.Context, request RegisterRequest, phone string) (UserResponse, error) {
	email := normalizeEmail(request.Email)
	name := strings.TrimSpace(request.FullName)
	if !validEmail(email) || name == "" || utf8.RuneCountInString(name) > 200 {
		return UserResponse{}, ErrInvalidInput
	}
	if !validPassword(request.Password) {
		return UserResponse{}, ErrWeakPassword
	}
	if s.repository == nil || s.tokens == nil {
		return UserResponse{}, ErrUnavailable
	}
	if s.email == nil || s.email.sender == nil {
		return UserResponse{}, ErrEmailNotConfigured
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("hash password: %w", err)
	}
	id, err := newID()
	if err != nil {
		return UserResponse{}, err
	}
	value := NewAccount{Account: Account{ID: id, Email: email, FullName: name, PasswordHash: string(hash), Role: RolePatient, Active: true}, Phone: phone}
	if phone != "" {
		patientID, err := newID()
		if err != nil {
			return UserResponse{}, err
		}
		value.PatientID = "BN-" + strings.ToUpper(strings.ReplaceAll(patientID, "-", ""))
	}
	verification, raw, err := s.email.prepare(id)
	if err != nil {
		return UserResponse{}, err
	}
	value.Verification = verification
	account, err := s.repository.Create(ctx, value)
	if err != nil {
		return UserResponse{}, err
	}
	if err := s.email.send(ctx, account, raw); err != nil {
		return UserResponse{}, err
	}
	return account.response(), nil
}

const dummyPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

func (s *Service) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	email := normalizeEmail(request.Email)
	if !validEmail(email) || request.Password == "" || len(request.Password) > MaxPasswordBytes {
		return LoginResponse{}, ErrInvalidLogin
	}
	if s.repository == nil || s.tokens == nil {
		return LoginResponse{}, ErrUnavailable
	}
	account, err := s.repository.FindByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(request.Password))
		return LoginResponse{}, ErrInvalidLogin
	}
	if err != nil {
		return LoginResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(request.Password)) != nil {
		return LoginResponse{}, ErrInvalidLogin
	}
	if !account.Active {
		return LoginResponse{}, ErrInactive
	}
	if !account.EmailVerified {
		return LoginResponse{}, ErrEmailUnverified
	}
	pair, err := s.tokens.issue(account)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{TokenResponse: pair, User: account.response()}, nil
}

func (s *Service) RefreshToken(ctx context.Context, request TokenRequest) (TokenResponse, error) {
	access, refresh, err := s.parsePair(request, false)
	if err != nil {
		return TokenResponse{}, err
	}
	account, err := s.activeAccount(ctx, refresh.Subject)
	if err != nil {
		return TokenResponse{}, err
	}
	for _, id := range []string{access.ID, refresh.ID} {
		revoked, err := s.repository.IsInvalidated(ctx, id)
		if err != nil {
			return TokenResponse{}, err
		}
		if revoked {
			return TokenResponse{}, ErrInvalidToken
		}
	}
	pair, err := s.tokens.issue(account)
	if err != nil {
		return TokenResponse{}, err
	}
	if err := s.repository.Revoke(ctx, []InvalidatedToken{refresh.invalidated(), access.invalidated()}, true); err != nil {
		return TokenResponse{}, err
	}
	return pair, nil
}

func (s *Service) Logout(ctx context.Context, request TokenRequest) error {
	access, refresh, err := s.parsePair(request, true)
	if err != nil {
		return err
	}
	if _, err := s.activeAccount(ctx, access.Subject); err != nil {
		return err
	}
	return s.repository.Revoke(ctx, []InvalidatedToken{refresh.invalidated(), access.invalidated()}, false)
}

func (s *Service) parsePair(request TokenRequest, allowExpiredRefresh bool) (*tokenClaims, *tokenClaims, error) {
	if request.AccessToken == "" || request.RefreshToken == "" {
		return nil, nil, ErrInvalidInput
	}
	if s.repository == nil || s.tokens == nil {
		return nil, nil, ErrUnavailable
	}
	access, err := s.tokens.parse(request.AccessToken, false, true)
	if err != nil {
		return nil, nil, err
	}
	refresh, err := s.tokens.parse(request.RefreshToken, true, allowExpiredRefresh)
	if err != nil {
		return nil, nil, err
	}
	if access.Subject != refresh.Subject {
		return nil, nil, ErrInvalidToken
	}
	return access, refresh, nil
}

func (s *Service) Verify(ctx context.Context, rawToken string) (middleware.Principal, error) {
	if s.repository == nil || s.tokens == nil {
		return middleware.Principal{}, ErrUnavailable
	}
	claims, err := s.tokens.parse(rawToken, false, false)
	if err != nil {
		return middleware.Principal{}, err
	}
	revoked, err := s.repository.IsInvalidated(ctx, claims.ID)
	if err != nil {
		return middleware.Principal{}, err
	}
	if revoked {
		return middleware.Principal{}, ErrInvalidToken
	}
	account, err := s.activeAccount(ctx, claims.Subject)
	if err != nil {
		return middleware.Principal{}, err
	}
	return middleware.Principal{UserID: account.ID, Roles: []string{string(account.Role)}}, nil
}

func (s *Service) activeAccount(ctx context.Context, id string) (Account, error) {
	account, err := s.repository.FindByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return Account{}, ErrInvalidToken
	}
	if err != nil {
		return Account{}, err
	}
	if !account.Active {
		return Account{}, ErrInactive
	}
	if !account.EmailVerified {
		return Account{}, ErrEmailUnverified
	}
	return account, nil
}

func normalizeEmail(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && len(value) <= 200
}
