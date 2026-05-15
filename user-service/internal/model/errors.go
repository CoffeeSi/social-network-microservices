package model

import "errors"

var (
	ErrFirstNameRequired = errors.New("first_name is required")
	ErrLastNameRequired  = errors.New("last_name is required")
	ErrEmailRequired     = errors.New("email is required")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrDOBRequired       = errors.New("dob is required")
	ErrTooYoung          = errors.New("user must be at least 13 years old")
)
