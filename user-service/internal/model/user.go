package model

import (
	"regexp"
	"time"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

const DateLayout = "02-01-2006"

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	DOB       time.Time
	IsActive  bool
	CreatedAt time.Time
}

func (u User) Validate() error {
	if u.FirstName == "" {
		return ErrFirstNameRequired
	}
	if u.LastName == "" {
		return ErrLastNameRequired
	}
	if u.Email == "" {
		return ErrEmailRequired
	}
	if !EmailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	if u.DOB.IsZero() {
		return ErrDOBRequired
	}
	minAge := time.Now().AddDate(-13, 0, 0)
	if u.DOB.After(minAge) {
		return ErrTooYoung
	}
	return nil
}
