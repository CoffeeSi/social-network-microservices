package config

import "os"

type Config struct {
	AuthGRPCPort string
	NatsURL      string
	SecretKey    string
	UserGRPCUrl  string
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
	return &Config{
		AuthGRPCPort: port,
		NatsURL:      natsURL,
		SecretKey:    secretKey,
		UserGRPCUrl:  userGRPCUrl,
	}
}
