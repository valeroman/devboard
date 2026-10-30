package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr string
}

func Load() (Config, error) {
	httpAddr := os.Getenv("IDENTITY_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8091"
	}

	if httpAddr[0] != ':' {
		return Config{}, fmt.Errorf("IDENTITY_HTTP_ADDR debe usar el formato de puertos :port")
	}

	return Config{
		HTTPAddr: httpAddr,
	}, nil
}
