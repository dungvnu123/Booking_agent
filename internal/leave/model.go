package leave

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Status string

const (
	Pending   Status = "PENDING"
	Approved  Status = "APPROVED"
	Rejected  Status = "REJECTED"
	Withdrawn Status = "WITHDRAWN"
)

type Request struct {
	ID       int64
	DoctorID string
	Status   Status
	FromDate pgtype.Date
	ToDate   pgtype.Date
}
type Item struct {
	ID                          int64
	LeaveRequestID              int64
	ProposedReplacementDoctorID *string
	ApprovedReplacementDoctorID *string
}
