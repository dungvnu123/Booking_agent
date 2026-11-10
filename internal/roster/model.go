package roster

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Status string

const (
	Pending   Status = "PENDING"
	Draft     Status = "DRAFT"
	Approved  Status = "APPROVED"
	Published Status = "PUBLISHED"
)

type Session string

const (
	Morning   Session = "SANG"
	Afternoon Session = "CHIEU"
)

type Availability string

const (
	Off       Availability = "OFF"
	OK        Availability = "OK"
	Preferred Availability = "PREF"
)

type MonthlyRoster struct {
	ID     int64
	Status Status
}
type WorkSchedule struct {
	ID       int64
	DoctorID string
	WorkDate pgtype.Date
	Status   string
}
type Slot struct {
	ID             int64
	WorkScheduleID int64
	Capacity       int
	BookedCount    int
}

// Do not enforce BookedCount <= Capacity: staff over-capacity booking is a supported use case.
