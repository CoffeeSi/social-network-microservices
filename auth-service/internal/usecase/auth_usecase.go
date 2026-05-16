package usecase

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/event"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/utils"
	"github.com/redis/go-redis/v9"
)

type AuthUsecase struct {
	publisher  EventPublisherInterface
	userClient UserClientInterface
	jwtToken   *utils.JwtToken
	redis      *redis.Client
	ttl        time.Duration
}

func NewAuthUsecase(publisher EventPublisherInterface, userClient UserClientInterface, jwtToken *utils.JwtToken, redis *redis.Client, ttl time.Duration) *AuthUsecase {
	return &AuthUsecase{publisher: publisher, userClient: userClient, jwtToken: jwtToken, redis: redis, ttl: ttl}
}
func generateVerificationCode(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
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

	userDOB, err := utils.ParseDOB(req.DOB)
	if err != nil {
		return err
	}

	newAuth := model.Auth{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       userDOB,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	err = uc.userClient.CreateUser(ctx, newAuth)
	if err != nil {
		return err
	}

	userRegisteredEvent := event.UserRegisteredEvent{
		EventType:  event.UserRegisteredEventType,
		OccurredAt: time.Now().Format(time.RFC3339),
		FirstName:  newAuth.FirstName,
		LastName:   newAuth.LastName,
		DOB:        newAuth.DOB.Format("2006-01-02"),
		Email:      newAuth.Email,
		Password:   newAuth.Password,
	}

	code, err := generateVerificationCode(6)
	if err != nil {
		return err
	}
	userVerificationEvent := event.UserVerificationEvent{
		EventType: event.UserVerificationEventType,
		OccuredAt: time.Now().Format(time.RFC3339),
		Email:     newAuth.Email,
		Code:      code,
	}
	key := fmt.Sprintf("verify:%s", newAuth.Email)
	codeMarshal, err := json.Marshal(code)
	if err != nil {
		log.Printf("error in marshaling")
	}
	uc.redis.Set(ctx, key, codeMarshal, uc.ttl)
	err = uc.publisher.Publish(event.UserRegisteredEventType, userRegisteredEvent)
	if err != nil {
		log.Printf("event publishing error: %v\n", err)
	}
	err = uc.publisher.Publish(event.UserVerificationEventType, userVerificationEvent)
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

	user, err := uc.userClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", "", err
	}
	if !user.IsActive {
		return "", "", errors.New("User is not active")
	}
	accessToken, err := uc.jwtToken.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.jwtToken.GenerateRefreshToken(user.ID)
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

func (uc *AuthUsecase) Verify(ctx context.Context, email, code string) (bool, error) {
	key := fmt.Sprintf("verify:%s", email)
	codeRedis, err := uc.redis.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("Error getting data from redis in verify method")
	}
	var res string
	json.Unmarshal(codeRedis, &res)
	if res == "" {
		return false, errors.New("expired code")
	} else if strings.Compare(res, code) == 0 {
		uc.userClient.ChangeStatus(ctx, email)
		uc.redis.Del(ctx, key)
		return true, nil
	} else {
		return false, errors.New("code does not match")
	}

}
