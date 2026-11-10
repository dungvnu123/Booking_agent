package patient

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

// CreatePatient is not implemented; port validation, authorization and transactions first.
func (s *Service) CreatePatient(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.CreatePatient: %w", fault.ErrNotImplemented)
}

// SearchPatients is not implemented; port validation, authorization and transactions first.
func (s *Service) SearchPatients(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.SearchPatients: %w", fault.ErrNotImplemented)
}

// GetPatientStats is not implemented; port validation, authorization and transactions first.
func (s *Service) GetPatientStats(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.GetPatientStats: %w", fault.ErrNotImplemented)
}

// GetAllPatients is not implemented; port validation, authorization and transactions first.
func (s *Service) GetAllPatients(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.GetAllPatients: %w", fault.ErrNotImplemented)
}

// GetPatientById is not implemented; port validation, authorization and transactions first.
func (s *Service) GetPatientById(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.GetPatientById: %w", fault.ErrNotImplemented)
}

// UpdatePatient is not implemented; port validation, authorization and transactions first.
func (s *Service) UpdatePatient(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.UpdatePatient: %w", fault.ErrNotImplemented)
}

// ImportPatients is not implemented; port validation, authorization and transactions first.
func (s *Service) ImportPatients(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.ImportPatients: %w", fault.ErrNotImplemented)
}

// ImportPatientsBatch is not implemented; port validation, authorization and transactions first.
func (s *Service) ImportPatientsBatch(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("patient.ImportPatientsBatch: %w", fault.ErrNotImplemented)
}
