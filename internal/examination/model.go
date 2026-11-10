package examination

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Examination struct {
	ID        int64
	PatientID string
	DoctorID  string
	CreatedAt time.Time
}

// TestResult is persisted in the existing table named test.
type TestResult struct {
	ID            int64
	ExaminationID int64
}
