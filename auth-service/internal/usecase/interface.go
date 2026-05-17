package usecase

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
)

type AuthUsecaseInterface interface {
	RegisterUser(ctx context.Context, req model.RegisterRequest) error
	LoginUser(ctx context.Context, req model.LoginRequest) (string, string, error)
	RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (string, error)
	Verify(ctx context.Context, email, code string) (bool, error)
}

type UserClientInterface interface {
	CreateUser(ctx context.Context, req model.Auth) error
	GetUserByEmail(ctx context.Context, email string) (model.Auth, error)
	ChangeStatus(ctx context.Context, email string) error
}

type EventPublisherInterface interface {
	Publish(subject string, payload any) error
}
