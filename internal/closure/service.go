package closure

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

// ListByMonth is not implemented; port validation, authorization and transactions first.
func (s *Service) ListByMonth(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("closure.ListByMonth: %w", fault.ErrNotImplemented)
}

// GetImpact is not implemented; port validation, authorization and transactions first.
func (s *Service) GetImpact(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("closure.GetImpact: %w", fault.ErrNotImplemented)
}

// GetImpactBulk is not implemented; port validation, authorization and transactions first.
func (s *Service) GetImpactBulk(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("closure.GetImpactBulk: %w", fault.ErrNotImplemented)
}

// Add is not implemented; port validation, authorization and transactions first.
func (s *Service) Add(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("closure.Add: %w", fault.ErrNotImplemented)
}

// AddBulk is not implemented; port validation, authorization and transactions first.
func (s *Service) AddBulk(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("closure.AddBulk: %w", fault.ErrNotImplemented)
}

// Delete is not implemented; port validation, authorization and transactions first.
func (s *Service) Delete(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("closure.Delete: %w", fault.ErrNotImplemented)
}
