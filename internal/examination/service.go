package examination

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
	return fmt.Errorf("examination.Create: %w", fault.ErrNotImplemented)
}

// GetByPatient is not implemented; port validation, authorization and transactions first.
func (s *Service) GetByPatient(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("examination.GetByPatient: %w", fault.ErrNotImplemented)
}

// GetById is not implemented; port validation, authorization and transactions first.
func (s *Service) GetById(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("examination.GetById: %w", fault.ErrNotImplemented)
}

// Update is not implemented; port validation, authorization and transactions first.
func (s *Service) Update(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("examination.Update: %w", fault.ErrNotImplemented)
}
