package newsletter

import "errors"

var (
	ErrEmptyEmail     = errors.New("email is empty")
	ErrorInvalidEmail = errors.New("invalid email")
)
