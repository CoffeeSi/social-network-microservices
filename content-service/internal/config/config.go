package config

import "os"

type Config struct {
	GRPCPort string
	MongoURI string
	DBName   string
}

func LoadConfig() Config {
	return Config{
		GRPCPort: getEnv("GRPC_PORT", "50053"),
		MongoURI: getEnv("MONGO_URI", "mongo://localhost:27017"),
		DBName:   getEnv("DB_NAME", "maxat_content_db"),
	}
}
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
