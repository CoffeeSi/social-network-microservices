package config

import "os"

type Config struct {
	HTTPPort           string
	AuthServiceAddr    string
	UserServiceAddr    string
	ContentServiceAddr string
	ChatServiceAddr    string
	SecretKey          string
}

func Load() Config {
	return Config{
		HTTPPort:           env("HTTP_PORT", "8080"),
		AuthServiceAddr:    env("AUTH_SERVICE_ADDR", "localhost:50051"),
		UserServiceAddr:    env("USER_SERVICE_ADDR", "localhost:50050"),
		ContentServiceAddr: env("CONTENT_SERVICE_ADDR", "localhost:50052"),
		ChatServiceAddr:    env("CHAT_SERVICE_ADDR", "localhost:50053"),
		SecretKey:          env("SECRET_KEY", "secret-key"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
