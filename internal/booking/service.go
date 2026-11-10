package booking

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

// GetAvailableSlots is not implemented; port validation, authorization and transactions first.
func (s *Service) GetAvailableSlots(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.GetAvailableSlots: %w", fault.ErrNotImplemented)
}

// GetSlotBoard is not implemented; port validation, authorization and transactions first.
func (s *Service) GetSlotBoard(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.GetSlotBoard: %w", fault.ErrNotImplemented)
}

// HoldSlot is not implemented; port validation, authorization and transactions first.
func (s *Service) HoldSlot(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.HoldSlot: %w", fault.ErrNotImplemented)
}

// ReleaseHold is not implemented; port validation, authorization and transactions first.
func (s *Service) ReleaseHold(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.ReleaseHold: %w", fault.ErrNotImplemented)
}

// CreateBooking is not implemented; port validation, authorization and transactions first.
func (s *Service) CreateBooking(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.CreateBooking: %w", fault.ErrNotImplemented)
}

// CreateStaffBooking is not implemented; port validation, authorization and transactions first.
func (s *Service) CreateStaffBooking(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.CreateStaffBooking: %w", fault.ErrNotImplemented)
}

// ListMyBookings is not implemented; port validation, authorization and transactions first.
func (s *Service) ListMyBookings(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.ListMyBookings: %w", fault.ErrNotImplemented)
}

// CancelBooking is not implemented; port validation, authorization and transactions first.
func (s *Service) CancelBooking(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.CancelBooking: %w", fault.ErrNotImplemented)
}

// ListBookingsForWorkSchedule is not implemented; port validation, authorization and transactions first.
func (s *Service) ListBookingsForWorkSchedule(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.ListBookingsForWorkSchedule: %w", fault.ErrNotImplemented)
}

// CheckIn is not implemented; port validation, authorization and transactions first.
func (s *Service) CheckIn(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.CheckIn: %w", fault.ErrNotImplemented)
}

// UndoCheckIn is not implemented; port validation, authorization and transactions first.
func (s *Service) UndoCheckIn(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.UndoCheckIn: %w", fault.ErrNotImplemented)
}

// MarkNoShow is not implemented; port validation, authorization and transactions first.
func (s *Service) MarkNoShow(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.MarkNoShow: %w", fault.ErrNotImplemented)
}

// CallBooking is not implemented; port validation, authorization and transactions first.
func (s *Service) CallBooking(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.CallBooking: %w", fault.ErrNotImplemented)
}

// FinishBooking is not implemented; port validation, authorization and transactions first.
func (s *Service) FinishBooking(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.FinishBooking: %w", fault.ErrNotImplemented)
}

// CallBookingFailed is not implemented; port validation, authorization and transactions first.
func (s *Service) CallBookingFailed(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.CallBookingFailed: %w", fault.ErrNotImplemented)
}

// StaffCancelBooking is not implemented; port validation, authorization and transactions first.
func (s *Service) StaffCancelBooking(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.StaffCancelBooking: %w", fault.ErrNotImplemented)
}

// GetAlternativeDoctors is not implemented; port validation, authorization and transactions first.
func (s *Service) GetAlternativeDoctors(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.GetAlternativeDoctors: %w", fault.ErrNotImplemented)
}

// GetAlternativeDoctorsForStaff is not implemented; port validation, authorization and transactions first.
func (s *Service) GetAlternativeDoctorsForStaff(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.GetAlternativeDoctorsForStaff: %w", fault.ErrNotImplemented)
}

// SwitchDoctor is not implemented; port validation, authorization and transactions first.
func (s *Service) SwitchDoctor(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.SwitchDoctor: %w", fault.ErrNotImplemented)
}

// StaffSwitchDoctor is not implemented; port validation, authorization and transactions first.
func (s *Service) StaffSwitchDoctor(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.StaffSwitchDoctor: %w", fault.ErrNotImplemented)
}

// Reschedule is not implemented; port validation, authorization and transactions first.
func (s *Service) Reschedule(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.Reschedule: %w", fault.ErrNotImplemented)
}

// StaffReschedule is not implemented; port validation, authorization and transactions first.
func (s *Service) StaffReschedule(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.StaffReschedule: %w", fault.ErrNotImplemented)
}

// GetQueue is not implemented; port validation, authorization and transactions first.
func (s *Service) GetQueue(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.GetQueue: %w", fault.ErrNotImplemented)
}

// ManualSwap is not implemented; port validation, authorization and transactions first.
func (s *Service) ManualSwap(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("booking.ManualSwap: %w", fault.ErrNotImplemented)
}
