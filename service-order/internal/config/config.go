package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Port string
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	cfg.App.Port = getEnv(
		"APP_PORT",
		"8080",
	)

	cfg.Database.Host = getEnv(
		"DB_HOST",
		"localhost",
	)

	cfg.Database.Port = getEnvAsInt(
		"DB_PORT",
		5432,
	)

	cfg.Database.User = getEnv(
		"DB_USER",
		"postgres",
	)

	cfg.Database.Password = getEnv(
		"DB_PASSWORD",
		"postgres",
	)

	cfg.Database.Name = getEnv(
		"DB_NAME",
		"order_db",
	)

	cfg.Database.SSLMode = getEnv(
		"DB_SSLMODE",
		"disable",
	)

	return &cfg, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(
	name string,
	defaultValue int,
) int {
	valueStr := os.Getenv(name)

	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
