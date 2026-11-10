package user

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

// GetRoles is not implemented; port validation, authorization and transactions first.
func (s *Service) GetRoles(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.GetRoles: %w", fault.ErrNotImplemented)
}

// GetUsers is not implemented; port validation, authorization and transactions first.
func (s *Service) GetUsers(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.GetUsers: %w", fault.ErrNotImplemented)
}

// SearchUsers is not implemented; port validation, authorization and transactions first.
func (s *Service) SearchUsers(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.SearchUsers: %w", fault.ErrNotImplemented)
}

// GetUser is not implemented; port validation, authorization and transactions first.
func (s *Service) GetUser(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.GetUser: %w", fault.ErrNotImplemented)
}

// CreateUser is not implemented; port validation, authorization and transactions first.
func (s *Service) CreateUser(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.CreateUser: %w", fault.ErrNotImplemented)
}

// UpdateUser is not implemented; port validation, authorization and transactions first.
func (s *Service) UpdateUser(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.UpdateUser: %w", fault.ErrNotImplemented)
}

// ChangeStatus is not implemented; port validation, authorization and transactions first.
func (s *Service) ChangeStatus(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.ChangeStatus: %w", fault.ErrNotImplemented)
}

// ResetPassword is not implemented; port validation, authorization and transactions first.
func (s *Service) ResetPassword(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("user.ResetPassword: %w", fault.ErrNotImplemented)
}
