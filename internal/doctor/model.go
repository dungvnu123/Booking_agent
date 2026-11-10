package doctor

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Profile struct {
	ID                          int64
	UserID                      string
	ConsultationDurationMinutes int
	PatientsPerHour             int
	IsActive                    bool
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}
