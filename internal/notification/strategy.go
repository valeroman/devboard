// Package notification strategy
package notification

import (
	"context"
	"log/slog"
)

// Notifier es la interfaz que define el contrato de la notification
type Notifier interface {
	Notify(ctx context.Context, userID string, message string) error
}

// LogNotifier registra las notificaciones en los logs
type LogNotifier struct {
	logger *slog.Logger
}

// NewLogNotifier constructor de LogNotifier
func NewLogNotifier(logger *slog.Logger) *LogNotifier {
	return &LogNotifier{logger: logger}
}

// Notify es el metodo para enviar notificaciones en los logs exportado
func (notifier *LogNotifier) Notify(ctx context.Context, userID string, message string) error {
	notifier.logger.InfoContext(ctx, "notificación enviada",
		slog.String("user_id", userID),
		slog.String("message", message),
	)

	return nil
}

// NoOpNotifier struct que no notifica pero ayuda a probar
type NoOpNotifier struct{}

// Notify metodo de NoOpNotifier
func (notifier *NoOpNotifier) Notify(_ context.Context, _ string, _ string) error {
	return nil
}
