package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/event"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/utils"
)

type AuthUsecase struct {
	publisher EventPublisherInterface
	jwtToken  *utils.JwtToken
}

func NewAuthUsecase(publisher EventPublisherInterface, jwtToken *utils.JwtToken) *AuthUsecase {
	return &AuthUsecase{publisher: publisher, jwtToken: jwtToken}
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, req model.RegisterRequest) error {
	if req.Email == "" ||
		req.Password == "" ||
		req.FirstName == "" ||
		req.LastName == "" ||
		req.DOB == "" {
		return errors.New("all fields are required")
	}
	if !utils.IsEmailValid(req.Email) {
		return errors.New("invalid email")
	}
	if !utils.IsPasswordValid(req.Password) {
		return errors.New("invalid password")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newAuth := model.Auth{
		Email:    req.Email,
		Password: hashedPassword,
	}

	userRegisteredEvent := event.UserRegisteredEvent{
		EventType:  event.UserRegisteredEventType,
		OccurredAt: time.Now().Format(time.RFC3339),
		Email:      newAuth.Email,
		Password:   newAuth.Password,
	}

	err = uc.publisher.Publish(event.UserRegisteredEventType, userRegisteredEvent)
	if err != nil {
		log.Printf("event publishing error: %v\n", err)
	}
	return nil
}

func (uc *AuthUsecase) LoginUser(ctx context.Context, req model.LoginRequest) (string, string, error) {
	if req.Email == "" || req.Password == "" {
		return "", "", errors.New("email and password are required")
	}
	if !utils.IsEmailValid(req.Email) {
		return "", "", errors.New("invalid email")
	}

	// TODO: fetch user from DB by email and verify password with utils.CheckPassword
	// Mocked userID for now
	userID := "mock-user-id"

	accessToken, err := uc.jwtToken.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.jwtToken.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (string, error) {
	userID, err := uc.jwtToken.VerifyToken(req.RefreshToken)
	if err != nil {
		return "", err
	}
	return uc.jwtToken.GenerateAccessToken(userID)
}
