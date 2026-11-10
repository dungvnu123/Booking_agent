package roster

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

// GetMyPick is not implemented; port validation, authorization and transactions first.
func (s *Service) GetMyPick(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetMyPick: %w", fault.ErrNotImplemented)
}

// SubmitPick is not implemented; port validation, authorization and transactions first.
func (s *Service) SubmitPick(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.SubmitPick: %w", fault.ErrNotImplemented)
}

// OpenMonth is not implemented; port validation, authorization and transactions first.
func (s *Service) OpenMonth(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.OpenMonth: %w", fault.ErrNotImplemented)
}

// ListMonths is not implemented; port validation, authorization and transactions first.
func (s *Service) ListMonths(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ListMonths: %w", fault.ErrNotImplemented)
}

// GetCurrentOpenMonth is not implemented; port validation, authorization and transactions first.
func (s *Service) GetCurrentOpenMonth(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetCurrentOpenMonth: %w", fault.ErrNotImplemented)
}

// GetMyPreview is not implemented; port validation, authorization and transactions first.
func (s *Service) GetMyPreview(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetMyPreview: %w", fault.ErrNotImplemented)
}

// GetMonth is not implemented; port validation, authorization and transactions first.
func (s *Service) GetMonth(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetMonth: %w", fault.ErrNotImplemented)
}

// GetAvgCapacity is not implemented; port validation, authorization and transactions first.
func (s *Service) GetAvgCapacity(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetAvgCapacity: %w", fault.ErrNotImplemented)
}

// UpdateSolverConfig is not implemented; port validation, authorization and transactions first.
func (s *Service) UpdateSolverConfig(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.UpdateSolverConfig: %w", fault.ErrNotImplemented)
}

// OverrideConfig is not implemented; port validation, authorization and transactions first.
func (s *Service) OverrideConfig(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.OverrideConfig: %w", fault.ErrNotImplemented)
}

// ClearOverrideConfig is not implemented; port validation, authorization and transactions first.
func (s *Service) ClearOverrideConfig(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ClearOverrideConfig: %w", fault.ErrNotImplemented)
}

// RunSolver is not implemented; port validation, authorization and transactions first.
func (s *Service) RunSolver(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.RunSolver: %w", fault.ErrNotImplemented)
}

// GetResult is not implemented; port validation, authorization and transactions first.
func (s *Service) GetResult(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetResult: %w", fault.ErrNotImplemented)
}

// GetPublishPreview is not implemented; port validation, authorization and transactions first.
func (s *Service) GetPublishPreview(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetPublishPreview: %w", fault.ErrNotImplemented)
}

// Approve is not implemented; port validation, authorization and transactions first.
func (s *Service) Approve(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.Approve: %w", fault.ErrNotImplemented)
}

// RevertToDraft is not implemented; port validation, authorization and transactions first.
func (s *Service) RevertToDraft(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.RevertToDraft: %w", fault.ErrNotImplemented)
}

// Publish is not implemented; port validation, authorization and transactions first.
func (s *Service) Publish(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.Publish: %w", fault.ErrNotImplemented)
}

// PublishNow is not implemented; port validation, authorization and transactions first.
func (s *Service) PublishNow(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.PublishNow: %w", fault.ErrNotImplemented)
}

// AddAssignment is not implemented; port validation, authorization and transactions first.
func (s *Service) AddAssignment(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.AddAssignment: %w", fault.ErrNotImplemented)
}

// RemoveAssignment is not implemented; port validation, authorization and transactions first.
func (s *Service) RemoveAssignment(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.RemoveAssignment: %w", fault.ErrNotImplemented)
}

// LockAssignment is not implemented; port validation, authorization and transactions first.
func (s *Service) LockAssignment(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.LockAssignment: %w", fault.ErrNotImplemented)
}

// LockCell is not implemented; port validation, authorization and transactions first.
func (s *Service) LockCell(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.LockCell: %w", fault.ErrNotImplemented)
}

// ListCellCandidates is not implemented; port validation, authorization and transactions first.
func (s *Service) ListCellCandidates(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ListCellCandidates: %w", fault.ErrNotImplemented)
}

// ListDoctorPickStatuses is not implemented; port validation, authorization and transactions first.
func (s *Service) ListDoctorPickStatuses(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ListDoctorPickStatuses: %w", fault.ErrNotImplemented)
}

// ReopenForDoctor is not implemented; port validation, authorization and transactions first.
func (s *Service) ReopenForDoctor(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ReopenForDoctor: %w", fault.ErrNotImplemented)
}

// LockForDoctor is not implemented; port validation, authorization and transactions first.
func (s *Service) LockForDoctor(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.LockForDoctor: %w", fault.ErrNotImplemented)
}

// UnlockForDoctor is not implemented; port validation, authorization and transactions first.
func (s *Service) UnlockForDoctor(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.UnlockForDoctor: %w", fault.ErrNotImplemented)
}

// GetPublishedSwapImpact is not implemented; port validation, authorization and transactions first.
func (s *Service) GetPublishedSwapImpact(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.GetPublishedSwapImpact: %w", fault.ErrNotImplemented)
}

// SwapPublishedDoctor is not implemented; port validation, authorization and transactions first.
func (s *Service) SwapPublishedDoctor(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.SwapPublishedDoctor: %w", fault.ErrNotImplemented)
}

// ListAuditLog is not implemented; port validation, authorization and transactions first.
func (s *Service) ListAuditLog(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ListAuditLog: %w", fault.ErrNotImplemented)
}

// ListWorkSchedules is not implemented; port validation, authorization and transactions first.
func (s *Service) ListWorkSchedules(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ListWorkSchedules: %w", fault.ErrNotImplemented)
}

// ListMyWorkSchedules is not implemented; port validation, authorization and transactions first.
func (s *Service) ListMyWorkSchedules(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("roster.ListMyWorkSchedules: %w", fault.ErrNotImplemented)
}
