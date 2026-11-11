package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"abclinic.local/backend/internal/platform/clock"
)

type EmailSender interface {
	Send(context.Context, string, string, string) error
}

type EmailConfig struct {
	VerificationURL string
	TokenTTL        time.Duration
	ResendCooldown  time.Duration
}

type EmailVerification struct {
	repository Repository
	sender     EmailSender
	clock      clock.Clock
	config     EmailConfig
}

func NewEmailVerification(repository Repository, sender EmailSender, c clock.Clock, config EmailConfig) (*EmailVerification, error) {
	link, err := url.Parse(config.VerificationURL)
	if err != nil || link.Host == "" || (link.Scheme != "https" && link.Scheme != "http") || link.User != nil || link.RawQuery != "" || link.Fragment != "" {
		return nil, errors.New("invalid email verification URL")
	}
	if repository == nil || c == nil || config.TokenTTL < time.Minute || config.TokenTTL > 24*time.Hour || config.ResendCooldown < time.Second || config.ResendCooldown >= config.TokenTTL {
		return nil, errors.New("email verification requires repository, clock and valid TTL/cooldown")
	}
	return &EmailVerification{repository: repository, sender: sender, clock: c, config: config}, nil
}

func (v *EmailVerification) prepare(accountID string) (EmailToken, string, error) {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return EmailToken{}, "", fmt.Errorf("generate verification token: %w", err)
	}
	raw := base64.RawURLEncoding.EncodeToString(bytes[:])
	hash, _ := emailTokenHash(raw)
	now := v.clock.Now()
	return EmailToken{AccountID: accountID, Hash: hash, ExpiresAt: now.Add(v.config.TokenTTL), SentAt: now}, raw, nil
}

func (v *EmailVerification) send(ctx context.Context, account Account, raw string) error {
	link, err := url.Parse(v.config.VerificationURL)
	if err != nil {
		return err
	}
	link.Fragment = "token=" + raw
	body := fmt.Sprintf("Xin chào,\n\nBạn đã yêu cầu xác thực email cho tài khoản Booking App.\nMở đường dẫn sau và chọn Xác nhận email:\n\n%s\n\nLink có hiệu lực %d phút và chỉ dùng được một lần.\nNếu bạn không đăng ký tài khoản, hãy bỏ qua email này.\n", link.String(), int(v.config.TokenTTL/time.Minute))
	if err := v.sender.Send(ctx, account.Email, "Xác thực email tài khoản Booking App", body); err != nil {
		slog.ErrorContext(ctx, "verification email delivery failed", "error", err)
		return ErrEmailDelivery
	}
	slog.InfoContext(ctx, "verification email accepted by SMTP")
	return nil
}

func emailTokenHash(raw string) (string, bool) {
	if len(raw) != 43 {
		return "", false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(decoded) != 32 {
		return "", false
	}
	hash := sha256.Sum256(decoded)
	return hex.EncodeToString(hash[:]), true
}

func (s *Service) VerifyEmail(ctx context.Context, request VerifyEmailRequest) error {
	hash, ok := emailTokenHash(request.Token)
	if !ok {
		return ErrInvalidEmailToken
	}
	if s.email == nil {
		return ErrEmailNotConfigured
	}
	return s.repository.VerifyEmail(ctx, hash, s.email.clock.Now())
}

func (s *Service) ResendVerification(ctx context.Context, request ResendVerificationRequest) error {
	email := normalizeEmail(request.Email)
	if !validEmail(email) {
		return ErrInvalidInput
	}
	if s.email == nil || s.email.sender == nil {
		return ErrEmailNotConfigured
	}
	account, err := s.repository.FindByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		slog.InfoContext(ctx, "verification email skipped", "reason", "account_not_found")
		return nil
	}
	if err != nil {
		return err
	}
	if account.EmailVerified {
		slog.InfoContext(ctx, "verification email skipped", "reason", "already_verified")
		return nil
	}
	if !account.Active {
		slog.InfoContext(ctx, "verification email skipped", "reason", "account_inactive")
		return nil
	}
	token, raw, err := s.email.prepare(account.ID)
	if err != nil {
		return err
	}
	stored, err := s.repository.ReplaceEmailToken(ctx, token, s.email.clock.Now().Add(-s.email.config.ResendCooldown))
	if err != nil {
		return err
	}
	if !stored {
		slog.InfoContext(ctx, "verification email skipped", "reason", "cooldown_or_account_changed")
		return nil
	}
	return s.email.send(ctx, account, raw)
}
