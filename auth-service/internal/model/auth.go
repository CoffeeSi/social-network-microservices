package model

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
	FirstName string
	LastName  string
	DOB       string
	Email     string
	Password  string
}
