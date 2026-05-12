package usecase

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/utils"
)

type AuthUsecase struct {
}

func NewAuthUsecase() *AuthUsecase {
	return &AuthUsecase{}
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, req model.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newAuth := model.Auth{
		Email:    req.Email,
		Password: hashedPassword,
	}
	_ = newAuth

	// TODO: send event to RabbitMQ to create new user

	return nil
}

func (uc *AuthUsecase) LoginUser(ctx context.Context, req model.LoginRequest) (string, error) {
	return "", nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (string, error) {
	return "", nil
}
