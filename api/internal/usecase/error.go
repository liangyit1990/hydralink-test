package usecase

import "github.com/pkg/errors"

// List of errors
var (
	ErrUserExists      = errors.New("user already exists")
	ErrEmailInvalid    = errors.New("invalid email")
	ErrPasswordInvalid = errors.New("invalid password")
)
