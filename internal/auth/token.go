package auth

import (
	"errors"
	"time"

	"abclinic.local/backend/internal/platform/clock"
	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	AccessSecret  string
	RefreshSecret string
	Issuer        string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Tokens struct {
	config TokenConfig
	clock  clock.Clock
}

type tokenClaims struct {
	Scope string `json:"scope,omitempty"`
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func NewTokens(config TokenConfig, c clock.Clock) (*Tokens, error) {
	if len(config.AccessSecret) < 32 || len(config.RefreshSecret) < 32 || config.AccessSecret == config.RefreshSecret ||
		config.Issuer == "" || config.AccessTTL < time.Second || config.RefreshTTL <= config.AccessTTL || c == nil {
		return nil, errors.New("JWT requires distinct secrets of at least 32 bytes, an issuer, clock and refresh TTL greater than access TTL")
	}
	return &Tokens{config: config, clock: c}, nil
}

func (t *Tokens) issue(account Account) (TokenResponse, error) {
	access, err := t.sign(account, false)
	if err != nil {
		return TokenResponse{}, err
	}
	refresh, err := t.sign(account, true)
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (t *Tokens) sign(account Account, refresh bool) (string, error) {
	id, err := newID()
	if err != nil {
		return "", err
	}
	secret, ttl := t.config.AccessSecret, t.config.AccessTTL
	if refresh {
		secret, ttl = t.config.RefreshSecret, t.config.RefreshTTL
	}
	now := t.clock.Now()
	claims := tokenClaims{RegisteredClaims: jwt.RegisteredClaims{
		ID: id, Subject: account.ID, Issuer: t.config.Issuer,
		IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}}
	if !refresh {
		claims.Scope = "ROLE_" + string(account.Role)
		claims.Email = account.Email
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func (t *Tokens) parse(raw string, refresh, allowExpired bool) (*tokenClaims, error) {
	secret := t.config.AccessSecret
	if refresh {
		secret = t.config.RefreshSecret
	}
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(*jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithoutClaimsValidation())
	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	now := t.clock.Now()
	if claims.Subject == "" || claims.ID == "" || claims.Issuer != t.config.Issuer || claims.IssuedAt == nil ||
		claims.ExpiresAt == nil || claims.IssuedAt.Time.After(now) || !claims.ExpiresAt.Time.After(claims.IssuedAt.Time) ||
		(claims.NotBefore != nil && claims.NotBefore.Time.After(now)) || (!allowExpired && !claims.ExpiresAt.Time.After(now)) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (c *tokenClaims) invalidated() InvalidatedToken {
	return InvalidatedToken{JTI: c.ID, ExpiresAt: c.ExpiresAt.Time}
}
