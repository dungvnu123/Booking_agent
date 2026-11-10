package attachment

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
// Go target design: direct ownership, independent of AiSession.
// These fields require an explicit migration before reusing existing attachments.
type Attachment struct {
	ID            int64
	PatientID     string
	ExaminationID *int64
	FileName      string
	StorageKey    string
	FileSize      int64
	ContentType   string
	CreatedAt     time.Time
}
