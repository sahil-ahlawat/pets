package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment        string
	Port               string
	DatabaseURL        string
	RedisAddr          string
	RedisPassword      string
	JWTSecret          string
	JWTAccessExpiryHr  int
	JWTRefreshExpiryHr int
	PublicBaseURL      string
	S3Bucket           string
	S3Endpoint         string
	S3Region           string
	S3AccessKey        string
	S3SecretKey        string
	AllowedOrigins     string
}

func Load() *Config {
	return &Config{
		Environment:        getEnv("ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/pet_platform?sslmode=disable"),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-pet-identity-platform-jwt-key-2026"),
		JWTAccessExpiryHr:  getEnvInt("JWT_ACCESS_EXPIRY_HR", 24),
		JWTRefreshExpiryHr: getEnvInt("JWT_REFRESH_EXPIRY_HR", 720),
		PublicBaseURL:      getEnv("PUBLIC_BASE_URL", "http://localhost:3000"),
		S3Bucket:           getEnv("S3_BUCKET", "pet-identity-assets"),
		S3Endpoint:         getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Region:           getEnv("S3_REGION", "us-east-1"),
		S3AccessKey:        getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey:        getEnv("S3_SECRET_KEY", "minioadmin"),
		AllowedOrigins:     getEnv("ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
