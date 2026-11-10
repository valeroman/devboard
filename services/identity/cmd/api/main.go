package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/valeroman/devboard/services/identity/internal/auth"
	"github.com/valeroman/devboard/services/identity/internal/cache"
	"github.com/valeroman/devboard/services/identity/internal/config"
	"github.com/valeroman/devboard/services/identity/internal/handler"
	"github.com/valeroman/devboard/services/identity/internal/middlewares"
	"github.com/valeroman/devboard/services/identity/internal/repository/postgres"
	"github.com/valeroman/devboard/services/identity/internal/server"
	"github.com/valeroman/devboard/services/identity/internal/usecase"
	"github.com/valeroman/devboard/services/identity/internal/validator"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

func main() {
	if err := run(); err != nil {
		slog.Error("Identity service se ha detenido", slog.Any("error", err))
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

	logger.Info("La conexion a la base de datos fue establecida")

	// Redis client
	redisClient, err := cache.NewClient(
		ctx,
		cache.Config{
			Addr:    cfg.RedisAddr,
			Timeout: 5 * time.Second,
		},
	)

	if err != nil {
		return err
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error(
				"no se pudo cerrar la conexión a Redis",
				slog.Any("error", err),
			)
		}
	}()

	logger.Info("conexion a Redis establecida")

	jwtManager, err := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:            []byte(cfg.JWTSecret),
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	})

	if err != nil {
		return err
	}

	tokenStore := auth.NewTokenStore(redisClient)

	// Repositories
	userRepository := postgres.NewUserRepository(pool)

	// Use cases
	userUseCase := usecase.NewUserUseCase(userRepository)

	authUseCase := usecase.NewAuthUseCase(
		userRepository,
		jwtManager,
		tokenStore,
		refreshTokenDuration,
	)

	validate := validator.New()

	// Handlers
	userHandler := handler.NewUserHandler(userUseCase, validate, logger)
	authHandler := handler.NewAuthHandler(authUseCase, validate, logger)

	// Server
	srv := server.New(cfg.HTTPAddr,
		server.Withlogger(logger),
		server.WithReadHeaderTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second))

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
				`{"status": "ok", "service": "identity"}`,
			))
		},
	))

	srv.RegisterRoutes("POST /api/v1/auth/register", http.HandlerFunc(authHandler.Register))
	srv.RegisterRoutes("POST /api/v1/auth/login", http.HandlerFunc(authHandler.Login))
	srv.RegisterRoutes("POST /api/v1/auth/refresh", http.HandlerFunc(authHandler.Refresh))
	srv.RegisterRoutes("POST /api/v1/auth/logout", http.HandlerFunc(authHandler.Logout))
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))

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

	logger.Info("Se inicio IDENTITY SERVICE", slog.String("addr", cfg.HTTPAddr))

	return srv.Start()

}
