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
	ID            string
	Email         string
	FullName      string
	Role          Role
	PasswordHash  string
	Active        bool
	EmailVerified bool
}

func (a Account) response() UserResponse {
	return UserResponse{UserID: a.ID, Email: a.Email, FullName: a.FullName, Role: a.Role, EmailVerified: a.EmailVerified}
}

type NewAccount struct {
	Account
	Phone        string
	PatientID    string
	Verification EmailToken
}

type InvalidatedToken struct {
	JTI       string
	ExpiresAt time.Time
}

type EmailToken struct {
	AccountID string
	Hash      string
	ExpiresAt time.Time
	SentAt    time.Time
}
