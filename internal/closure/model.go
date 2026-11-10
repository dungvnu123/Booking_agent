package closure

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Session string

const (
	FullDay   Session = "FULL_DAY"
	Morning   Session = "SANG"
	Afternoon Session = "CHIEU"
)

type Closure struct {
	ID        int64
	Date      pgtype.Date
	Session   Session
	CreatedBy string
	Reason    string
}
