package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	RedisAddr   string
	JWTSecret   string
}

func Load() (Config, error) {

	cfg := Config{
		HTTPAddr:    envOrDeffault("IDENTITY_HTTP_ADDR", ":8091"),
		DatabaseURL: envOrDeffault("IDENTITY_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/identity?sslmode=disable"),
		RedisAddr:   envOrDeffault("REDIS_ADDR", "localhost:6379"),
		JWTSecret:   envOrDeffault("JWT_SECRET", "dev-secret-key-change-in-production-min-32-chars$$"),
	}

	// httpAddr := os.Getenv("IDENTITY_HTTP_ADDR")
	// if httpAddr == "" {
	// 	httpAddr = ":8091"
	// }

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("IDENTITY_DATABASE_URL es requerida")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("IDENTITY_JWT_SECRET es requerida")
	}

	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("IDENTITY_JWT_SECRET debe tener al menos 32 caracteres")
	}

	if cfg.HTTPAddr[0] != ':' {
		return Config{}, fmt.Errorf("IDENTITY_HTTP_ADDR debe usar el formato de puertos :port")
	}

	return cfg, nil
}

func envOrDeffault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
