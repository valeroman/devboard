// Package main api
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/valeroman/devboard/docs"
	"github.com/valeroman/devboard/internal/handler"
	"github.com/valeroman/devboard/internal/logger"
	"github.com/valeroman/devboard/internal/middlewares"
	"github.com/valeroman/devboard/internal/repository/memory"
	"github.com/valeroman/devboard/internal/server"
	"github.com/valeroman/devboard/internal/usecase"
	"github.com/valeroman/devboard/internal/validator"
)

// @title Devboard API
// @version 1.0
// @description API REST para la gestión de proyectos
// @contact.me Soporte Devboard RV
// @host localhost:8090
// @BasePath /api/v1
func main() {

	// 1. Configuración del entorno
	// 2. Logger
	log := setupLogger()

	// 3. Dependencies compartidas

	// Infrastructure: adapters
	userRepo := memory.NewUserRepository()

	// Use cases: lógica del negocio
	userUC := usecase.NewUserUseCase(userRepo)

	// Validación compartida
	validate := validator.New()

	// 4. Handlers
	healthHandler := handler.NewHealtHandler()
	userHandler := handler.NewUserHandler(userUC, validate, log)

	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	// 5. Servidor con Functional options
	// Aqui arranca el servidor
	srv := server.New(":8090",
		server.Withlogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	// 6. Registro de rutas
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))

	// 7. Middleware chain
	srv.UseChain(
		middlewares.Recovery(log),
		middlewares.Logger(log),
	)

	// 8. Arrancar shutdown limpio
	if err := srv.Start(); err != nil {
		log.Error("error fatal", "error", err)
		os.Exit(1)
	}

}

func setupLogger() *slog.Logger {
	if os.Getenv("GO_ENV") == "production" {
		return logger.New(logger.ProductionConfig())
	}

	return logger.New(logger.DefaultConfig())
}
