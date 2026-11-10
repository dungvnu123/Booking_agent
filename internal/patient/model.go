package patient

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
// ID is a business string such as BN-XXXXXXXX, not necessarily a UUID.
type Patient struct {
	ID       string
	UserID   *string
	FullName *string
	Phone    *string
	// DateOfBirth is date-only in the clinic timezone, not a timestamp.
	DateOfBirth *pgtype.Date
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
