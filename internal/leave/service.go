package leave

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

// Create is not implemented; port validation, authorization and transactions first.
func (s *Service) Create(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.Create: %w", fault.ErrNotImplemented)
}

// ListMine is not implemented; port validation, authorization and transactions first.
func (s *Service) ListMine(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.ListMine: %w", fault.ErrNotImplemented)
}

// ListAll is not implemented; port validation, authorization and transactions first.
func (s *Service) ListAll(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.ListAll: %w", fault.ErrNotImplemented)
}

// Preview is not implemented; port validation, authorization and transactions first.
func (s *Service) Preview(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.Preview: %w", fault.ErrNotImplemented)
}

// Impact is not implemented; port validation, authorization and transactions first.
func (s *Service) Impact(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.Impact: %w", fault.ErrNotImplemented)
}

// Approve is not implemented; port validation, authorization and transactions first.
func (s *Service) Approve(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.Approve: %w", fault.ErrNotImplemented)
}

// Reject is not implemented; port validation, authorization and transactions first.
func (s *Service) Reject(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.Reject: %w", fault.ErrNotImplemented)
}

// Withdraw is not implemented; port validation, authorization and transactions first.
func (s *Service) Withdraw(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("leave.Withdraw: %w", fault.ErrNotImplemented)
}
