package repository

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidID        = errors.New("invalid user id")
	ErrDuplicateEmail   = errors.New("email already exists")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidDOB       = errors.New("invalid date of birth")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
	ErrOnUpdate         = errors.New("failed to update")
)
