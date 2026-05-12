package model

type RegisterRequest struct {
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type Auth struct {
	Email    string
	Password string
}
