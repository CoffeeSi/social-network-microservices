package model

import (
	"errors"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

const DateLayout = "02-01-2006"

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	DOB       string
	IsActive  bool
}

func (u User) Validate() error {
	if u.FirstName == "" {
		return errors.New("first_name is required")
	}
	if u.LastName == "" {
		return errors.New("last_name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	if u.DOB == "" {
		return errors.New("dob is required")
	}
	t, err := time.Parse(DateLayout, u.DOB)
	if err != nil {
		return errors.New("invalid dob format, expected DD-MM-YYYY")
	}
	minAge := time.Now().AddDate(-13, 0, 0)
	if t.After(minAge) {
		return errors.New("user must be at least 13 years old")
	}
	return nil
}
