package auth

import "time"

type Role string

const (
	RoleAdmin        Role = "ADMIN"
	RoleManager      Role = "MANAGER"
	RoleDoctor       Role = "DOCTOR"
	RoleReceptionist Role = "RECEPTIONIST"
	RolePatient      Role = "PATIENT"
)

type Account struct {
	ID           string
	Email        string
	FullName     string
	Role         Role
	PasswordHash string
	Active       bool
}

func (a Account) response() UserResponse {
	return UserResponse{UserID: a.ID, Email: a.Email, FullName: a.FullName, Role: a.Role}
}

type NewAccount struct {
	Account
	Phone     string
	PatientID string
}

type InvalidatedToken struct {
	JTI       string
	ExpiresAt time.Time
}
