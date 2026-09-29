// Package logger provides structured logging helpers for the application.
package logger

import (
	"io"
	"log/slog"
	"os"
)

// Config contains the logger configuration.
type Config struct {
	Level  slog.Level
	JSON   bool
	Output io.Writer
}

// DefaultConfig returns a default development logger configuration.
func DefaultConfig() Config {
	return Config{
		Level:  slog.LevelDebug,
		JSON:   false,
		Output: os.Stdout,
	}
}

// ProductionConfig returns a default production logger configuration.
func ProductionConfig() Config {
	return Config{
		Level:  slog.LevelInfo,
		JSON:   true,
		Output: os.Stdout,
	}
}

// New creates a new slog logger using the provided configuration.
func New(config Config) *slog.Logger {
	out := config.Output
	if out == nil {
		out = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.Level == slog.LevelDebug,
	}

	var handler slog.Handler
	if config.JSON {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = slog.NewTextHandler(out, opts)
	}

	return slog.New(handler)
}
