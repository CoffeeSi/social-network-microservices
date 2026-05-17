package config

import (
	"os"
	"time"
)

type Config struct {
	GRPCPort    string
	NatsURL     string
	SecretKey   string
	UserGRPCUrl string
	REDIS       string
	TTL         time.Duration
}

func NewConfig() *Config {
	port := os.Getenv("AUTH_GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "secret-key"
	}
	userGRPCUrl := os.Getenv("USER_GRPC_URL")
	if userGRPCUrl == "" {
		userGRPCUrl = "localhost:50052"
	}
	REDIS := os.Getenv("REDIS")
	if REDIS == "" {
		REDIS = "redis://redis:6379"
	}
	TTL := os.Getenv("TTL")
	if TTL == "" {
		TTL = "10m"
	}
	duration, _ := time.ParseDuration(TTL)
	return &Config{
		GRPCPort:    port,
		NatsURL:     natsURL,
		SecretKey:   secretKey,
		UserGRPCUrl: userGRPCUrl,
		REDIS:       REDIS,
		TTL:         duration,
	}
}
