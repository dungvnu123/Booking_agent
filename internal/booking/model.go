package booking

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Status string

const (
	Booked    Status = "BOOKED"
	CheckedIn Status = "CHECKED_IN"
	Late      Status = "LATE"
	Called    Status = "CALLED"
	Completed Status = "COMPLETED"
	NoShow    Status = "NO_SHOW"
	Cancelled Status = "CANCELLED"
)

type Booking struct {
	ID             int64
	PatientID      string
	SlotID         int64
	Status         Status
	CreatedBy      *string
	OriginalSlotID *int64
	CheckedInAt    *time.Time
	QueueOrderAt   *time.Time
	CheckedInLate  *bool
}
type SlotHold struct {
	ID        int64
	SlotID    int64
	PatientID string
	ExpiresAt time.Time
}
