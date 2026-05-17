package app

import "os"

type Config struct {
	Port      string
	MongoURI  string
	Migration string
	NatsURI   string
	REDIS     string
	TTL       string
	RPM       string
}

func NewConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", ":50050"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		NatsURI:  getEnv("NATS_URI", "nats://localhost:4222"),
		REDIS:    getEnv("REDIS", "redis://localhost:6379"),
		TTL:      getEnv("CACHE_TTL_SECONDS", "60"),
		RPM:      getEnv("RATE_LIMIT_RPM", "100"),
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
