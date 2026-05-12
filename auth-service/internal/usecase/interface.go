package usecase

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
)

type AuthUsecaseInterface interface {
	RegisterUser(ctx context.Context, req model.RegisterRequest) error
	LoginUser(ctx context.Context, req model.LoginRequest) (string, error)
	RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (string, error)
}
