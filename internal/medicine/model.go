package medicine

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Type string

const (
	Medicine   Type = "MEDICINE"
	Supplement Type = "SUPPLEMENT"
)

// Numeric preserves decimal precision. JSON DTOs must map to the agreed
// frontend representation; never convert medicine prices to float64.
type MedicineRecord struct {
	ID           int64
	Code         string
	Type         Type
	ImportPrice  pgtype.Numeric
	SellingPrice pgtype.Numeric
}
