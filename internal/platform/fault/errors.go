package fault

import "errors"

var ErrNotImplemented = errors.New("use case is not implemented")

type Error struct {
	Status  int
	Code    int
	Message string
}

func (e *Error) Error() string { return e.Message }
