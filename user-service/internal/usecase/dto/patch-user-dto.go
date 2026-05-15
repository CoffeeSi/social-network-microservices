package dto

import "time"

type PatchUserDTO struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	DOB       time.Time
}
