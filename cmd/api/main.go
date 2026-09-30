// Package main api
package main

import (
	"log/slog"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/valeroman/devboard/docs"
	"github.com/valeroman/devboard/internal/handler"
	"github.com/valeroman/devboard/internal/logger"
	"github.com/valeroman/devboard/internal/middlewares"
	"github.com/valeroman/devboard/internal/server"
)

// @title Devboard API
// @version 1.0
// @description API REST para la gestión de proyectos
// @contact.me Soporte Devboard RV
// @host localhost:8090
// @BasePath /api/v1
func main() {

	// 1. Configuración del entorno
	env := os.Getenv("GO_ENV")
	slog.Info("entorno cargado", "GO_ENV", env)

	// 2. Logger
	var log *slog.Logger
	if env == "production" {
		log = logger.New(logger.ProductionConfig())
	} else {
		log = logger.New(logger.DefaultConfig())
	}

	// 3. Dependencies compartidas
	// validate := validator.New()
	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	// 4. Servidor con Functional options
	// Aqui arranca el servidor
	srv := server.New(":8090",
		server.Withlogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	// 5. Handlers
	healthHandler := handler.NewHealtHandler()

	// 6. Registro de rutas
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

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
