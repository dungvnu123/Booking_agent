package doctor

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

// ListDoctorProfiles is not implemented; port validation, authorization and transactions first.
func (s *Service) ListDoctorProfiles(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("doctor.ListDoctorProfiles: %w", fault.ErrNotImplemented)
}

// AssignProfile is not implemented; port validation, authorization and transactions first.
func (s *Service) AssignProfile(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("doctor.AssignProfile: %w", fault.ErrNotImplemented)
}
