package auth

import (
	"context"
	"errors"

	"abclinic.local/backend/internal/platform/fault"
)

var (
	ErrInvalidInput    = &fault.Error{Status: 400, Code: 400, Message: "Invalid registration or login data"}
	ErrInvalidLogin    = &fault.Error{Status: 401, Code: 2002, Message: "Email or password is incorrect"}
	ErrInvalidToken    = &fault.Error{Status: 401, Code: 2000, Message: "Invalid or expired token"}
	ErrInactive        = &fault.Error{Status: 403, Code: 2004, Message: "Account is inactive"}
	ErrEmailExists     = &fault.Error{Status: 409, Code: 3000, Message: "Email already exists"}
	ErrPhoneExists     = &fault.Error{Status: 409, Code: 3003, Message: "Phone number already exists"}
	ErrPatientConflict = &fault.Error{Status: 409, Code: 3003, Message: "Phone number matches multiple patient records; contact the clinic"}
	ErrRoleMissing     = &fault.Error{Status: 503, Code: 3002, Message: "PATIENT role is not configured"}
	ErrUnavailable     = &fault.Error{Status: 503, Code: 1503, Message: "Authentication is not configured or database is unavailable"}
	ErrNotFound        = errors.New("account not found")
)

type Repository interface {
	Create(context.Context, NewAccount) (Account, error)
	FindByEmail(context.Context, string) (Account, error)
	FindByID(context.Context, string) (Account, error)
	IsInvalidated(context.Context, string) (bool, error)
	Revoke(context.Context, []InvalidatedToken, bool) error
}
