package user

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
// User IDs are UUID strings; never substitute the demo's numeric IDs.
type User struct {
	ID           string
	RoleID       *int64
	Username     *string
	FullName     *string
	Email        *string
	PasswordHash string
	Phone        *string
	IsActive     bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
