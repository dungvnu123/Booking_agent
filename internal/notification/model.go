package notification

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Notification struct {
	ID               int64
	RecipientID      string
	Type             string
	Title            string
	Message          string
	RelatedBookingID *int64
	IsRead           bool
	CreatedAt        time.Time
}
