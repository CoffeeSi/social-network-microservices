package repository

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidID      = errors.New("invalid user id")
	ErrDuplicateEmail = errors.New("email already exists")
)
