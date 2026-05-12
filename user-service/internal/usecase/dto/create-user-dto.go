package dto

import "time"

type CreateUserDTO struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	DOB       time.Time
	IsActive  bool
	CreatedAt time.Time
}
