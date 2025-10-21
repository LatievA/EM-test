package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	DBName           string
	SSLMode          string
	ConnectionString string
}

func Load() (*PostgresConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &PostgresConfig{
		Host:     getEnv("PG_HOST", "localhost"),
		Port:     getEnvAsInt("PG_PORT", 5432),
		User:     getEnv("PG_USER", "postgres"),
		Password: getEnv("PG_PASSWORD", "postgres"),
		DBName:   getEnv("PG_DATABASE", "postgres"),
		SSLMode:  getEnv("PG_SSLMODE", "disable"),
		ConnectionString: getEnv("DATABASE_URL", ""),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if vStr, ok := os.LookupEnv(key); ok && vStr != "" {
		if v, err := strconv.Atoi(vStr); err == nil {
			return v
		}
	}
	return fallback
}
