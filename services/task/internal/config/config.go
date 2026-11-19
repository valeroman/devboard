package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddr             string
	DatabaseURL          string
	JWTSecret            string
	WorkspaceGRPCTarget  string
	WorkspaceGRPCTimeout time.Duration
}

func Load() (Config, error) {

	workspaceGRPCTimeoutValue := envOrDeffault("WORKSPACE_GRPC_TIMEOUT", "1s")
	workspaceGRPCTimerout, err := time.ParseDuration(workspaceGRPCTimeoutValue)

	if err != nil {
		return Config{}, fmt.Errorf("parse WORKSPACE_GRPC_TIMEOUT: %w", err)
	}

	cfg := Config{
		HTTPAddr:             envOrDeffault("TASK_HTTP_ADDR", ":8093"),
		DatabaseURL:          envOrDeffault("TASK_DATABASE_URL", ""),
		JWTSecret:            envOrDeffault("JWT_SECRET", "dev-secret-key-change-in-production-min-32-chars$$"),
		WorkspaceGRPCTarget:  envOrDeffault("WORKSPACE_GRPC_TARGET", "localhost:9072"),
		WorkspaceGRPCTimeout: workspaceGRPCTimerout,
	}

	// httpAddr := os.Getenv("IDENTITY_HTTP_ADDR")
	// if httpAddr == "" {
	// 	httpAddr = ":8091"
	// }

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("TASK_DATABASE_URL es requerida")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET es requerida")
	}

	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET debe tener al menos 32 caracteres")
	}

	if cfg.HTTPAddr[0] != ':' {
		return Config{}, fmt.Errorf("TASK_HTTP_ADDR debe usar el formato de puertos :port")
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
