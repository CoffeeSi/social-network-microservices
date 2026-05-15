package usecase

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrTooYoung        = errors.New("user must be at least 13 years old")
	ErrIDEmpty         = errors.New("id is empty")
	ErrInvalidEmail    = errors.New("invalid email")
)
