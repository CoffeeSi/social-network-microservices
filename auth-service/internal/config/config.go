package config

import "os"

type Config struct {
	GRPCPort  string
	NatsURL   string
	SecretKey string
}

func NewConfig() *Config {
	port := os.Getenv("GRPC_PORT")
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
	return &Config{
		GRPCPort:  port,
		NatsURL:   natsURL,
		SecretKey: secretKey,
	}
}
