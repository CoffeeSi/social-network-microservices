package config

import "os"

type Config struct {
	GRPCPort    string
	MongoURI    string
	MongoName   string
	SecretKey   string
	UserGRPCUrl string
}

func NewConfig() *Config {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoName := os.Getenv("MONGO_NAME")
	if mongoName == "" {
		mongoName = "social_network"
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
		GRPCPort:    port,
		MongoURI:    mongoURI,
		MongoName:   mongoName,
		SecretKey:   secretKey,
		UserGRPCUrl: userGRPCUrl,
	}
}
