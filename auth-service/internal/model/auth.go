package model

import "time"

type RegisterRequest struct {
	FirstName string
	LastName  string
	DOB       string
	Email     string
	Password  string
}

type LoginRequest struct {
	Email    string
	Password string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type Auth struct {
	ID        string
	FirstName string
	LastName  string
	DOB       time.Time
	Email     string
	Password  string
	IsActive  bool
}
