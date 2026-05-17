package app

import "os"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	NatsURI  string
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getEnv("SMTP_PORT", "587"),
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", ""),
		NatsURI:  getEnv("NATS_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
