package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

func queryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("empty date")
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", value)
}

func authContext(c *gin.Context) context.Context {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return c.Request.Context()
	}
	return metadata.AppendToOutgoingContext(c.Request.Context(), "authorization", auth)
}

func (g *Gateway) requireUserID(c *gin.Context) (string, bool) {
	userID, err := g.userIDFromRequest(c)
	if err != nil {
		writeError(c, http.StatusUnauthorized, err.Error())
		return "", false
	}
	return userID, true
}

func (g *Gateway) userIDFromRequest(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if tokenString == "" {
		return "", errors.New("authorization token is required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(g.secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid authorization token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("user id not found in token")
	}
	return userID, nil
}
