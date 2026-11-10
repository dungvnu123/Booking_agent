package auth

import (
	"time"
)

// Initial domain outline from the supplied inventory, not a full DB mapping.
// Keep persistence models separate from JSON request/response DTOs.
type Role string

const (
	RoleAdmin        Role = "ADMIN"
	RoleManager      Role = "MANAGER"
	RoleDoctor       Role = "DOCTOR"
	RoleReceptionist Role = "RECEPTIONIST"
	RolePatient      Role = "PATIENT"
)

type InvalidatedToken struct {
	JTI       string
	ExpiresAt time.Time
}
