package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtToken struct {
	secretKey []byte
}

func NewJWTToken(secretKey string) *JwtToken {
	return &JwtToken{
		secretKey: []byte(secretKey),
	}
}

func (t *JwtToken) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  userID,
			"type": "access",
			"exp":  time.Now().Add(15 * time.Minute).Unix(),
		})
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *JwtToken) GenerateRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  userID,
			"type": "refresh",
			"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		})
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *JwtToken) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{},
		func(token *jwt.Token) (any, error) {
			return t.secretKey, nil
		})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return userID, nil
}
