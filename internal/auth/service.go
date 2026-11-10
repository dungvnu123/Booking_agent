package auth

import (
	"abclinic.local/backend/internal/platform/clock"
	"abclinic.local/backend/internal/platform/fault"
	"context"
	"fmt"
)

// Service owns this module's business use cases. Add repository interfaces only
// when the first use case is implemented; no permissive fake data is returned.
type Service struct{ clock clock.Clock }

func NewService(c clock.Clock) *Service { return &Service{clock: c} }

// Register is not implemented; port validation, authorization and transactions first.
func (s *Service) Register(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("auth.Register: %w", fault.ErrNotImplemented)
}

// RegisterPatient is not implemented; port validation, authorization and transactions first.
func (s *Service) RegisterPatient(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("auth.RegisterPatient: %w", fault.ErrNotImplemented)
}

// Login is not implemented; port validation, authorization and transactions first.
func (s *Service) Login(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("auth.Login: %w", fault.ErrNotImplemented)
}

// RefreshToken is not implemented; port validation, authorization and transactions first.
func (s *Service) RefreshToken(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("auth.RefreshToken: %w", fault.ErrNotImplemented)
}

// Logout is not implemented; port validation, authorization and transactions first.
func (s *Service) Logout(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("auth.Logout: %w", fault.ErrNotImplemented)
}
