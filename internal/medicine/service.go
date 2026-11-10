package medicine

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

// GetMedicines is not implemented; port validation, authorization and transactions first.
func (s *Service) GetMedicines(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.GetMedicines: %w", fault.ErrNotImplemented)
}

// GetMetadata is not implemented; port validation, authorization and transactions first.
func (s *Service) GetMetadata(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.GetMetadata: %w", fault.ErrNotImplemented)
}

// GetMedicineById is not implemented; port validation, authorization and transactions first.
func (s *Service) GetMedicineById(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.GetMedicineById: %w", fault.ErrNotImplemented)
}

// CreateMedicine is not implemented; port validation, authorization and transactions first.
func (s *Service) CreateMedicine(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.CreateMedicine: %w", fault.ErrNotImplemented)
}

// UpdateMedicine is not implemented; port validation, authorization and transactions first.
func (s *Service) UpdateMedicine(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.UpdateMedicine: %w", fault.ErrNotImplemented)
}

// DeleteMedicine is not implemented; port validation, authorization and transactions first.
func (s *Service) DeleteMedicine(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.DeleteMedicine: %w", fault.ErrNotImplemented)
}

// BulkImport is not implemented; port validation, authorization and transactions first.
func (s *Service) BulkImport(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("medicine.BulkImport: %w", fault.ErrNotImplemented)
}
