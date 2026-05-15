package model

import "time"

type User struct {
	ID        string
	FirstName string
	LastName  string
	DOB       time.Time
	Email     string
	Password  string
}
