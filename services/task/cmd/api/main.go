package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/valeroman/devboard/services/task/internal/auth"
	"github.com/valeroman/devboard/services/task/internal/config"
	"github.com/valeroman/devboard/services/task/internal/handler"
	"github.com/valeroman/devboard/services/task/internal/middlewares"
	"github.com/valeroman/devboard/services/task/internal/repository/postgres"
	"github.com/valeroman/devboard/services/task/internal/server"
	"github.com/valeroman/devboard/services/task/internal/usecase"
	"github.com/valeroman/devboard/services/task/internal/validator"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Task service se ha detenido", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()

	pool, err := postgres.NewPool(
		ctx,
		postgres.Config{
			URL:             cfg.DatabaseURL,
			MaxConns:        10,
			MinConns:        2,
			MaxConnLifeTime: time.Hour,
			MaxConnIdleTime: 30 * time.Minute,
		},
	)

	if err != nil {
		return err
	}

	defer pool.Close()

	logger.Info("La conexion a la base de datos de task fue establecida")

	tokenVerifier, err := auth.NewTokenVerifier([]byte(cfg.JWTSecret))

	if err != nil {
		return err
	}

	// Repositories
	taskRepository := postgres.NewTaskRepository(pool)
	commentRepository := postgres.NewCommentRepository(pool)

	// Use cases
	taskUseCase := usecase.NewTaskUseCase(taskRepository, taskRepository)
	commenttUseCase := usecase.NewCommentUseCase(commentRepository, taskRepository)

	validate := validator.New()

	// Handlers
	taskHandler := handler.NewTaskHandler(taskUseCase, validate, logger)
	commentHandler := handler.NewCommentHandler(commenttUseCase, validate, logger)

	// Server
	srv := server.New(cfg.HTTPAddr,
		server.Withlogger(logger),
		server.WithReadHeaderTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second))

	authMiddleware := middlewares.RequireAuth(tokenVerifier, logger)

	// Rutas Publicas
	srv.RegisterRoutes("GET /health", http.HandlerFunc(
		func(write http.ResponseWriter, _ *http.Request) {
			write.Header().Set(
				"Content-Type", "application/json",
			)

			write.WriteHeader(
				http.StatusOK,
			)

			_, _ = write.Write([]byte(
				`{"status": "ok", "service": "task"}`,
			))
		},
	))

	// Task routes
	srv.RegisterRoutes("POST /api/v1/tasks", middlewares.Chain(
		http.HandlerFunc(taskHandler.Create),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/tasks/{id}", middlewares.Chain(
		http.HandlerFunc(taskHandler.Get),
		authMiddleware,
	))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/status", middlewares.Chain(
		http.HandlerFunc(taskHandler.UpdateStatus),
		authMiddleware,
	))

	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/assign", middlewares.Chain(
		http.HandlerFunc(taskHandler.Assign),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/projects/{id}/tasks", middlewares.Chain(
		http.HandlerFunc(taskHandler.ListByProject),
		authMiddleware,
	))

	// Comments routes
	srv.RegisterRoutes("POST /api/v1/tasks/{id}/comments", middlewares.Chain(
		http.HandlerFunc(commentHandler.Create),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/tasks/{id}/comments", middlewares.Chain(
		http.HandlerFunc(commentHandler.ListByTask),
		authMiddleware,
	))

	srv.RegisterRoutes("DELETE /api/v1/comments/{id}", middlewares.Chain(
		http.HandlerFunc(commentHandler.Delete),
		authMiddleware,
	))

	// HTTP middleware
	srv.Use(
		middlewares.SecurityHeaders(middlewares.SecurityConfig{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:5173",
			},
		}),
	)

	srv.Use(
		middlewares.Logger(logger),
	)

	srv.Use(
		middlewares.Recovery(logger),
	)

	logger.Info("Se inicio TASK SERVICE", slog.String("addr", cfg.HTTPAddr))

	return srv.Start()

}
