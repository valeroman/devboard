// Package main api
package main

import (
	"os"

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

	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelInfo,
	// }))

	logger := logger.New(logger.DefaultConfig())

	// Aqui arranca el servidor
	srv := server.New(":8090", logger)

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	healthHander := handler.NewHealtHandler()

	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHander)
	srv.RegisterRoutes("GET /ready", healthHander)

	// Ruta temporal para probar Recovery
	// srv.RegisterRoutes("GET /panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	panic("error provocado para probar Recovery")
	// }))

	if err := srv.Start(); err != nil {
		logger.Error("error fatal", "error", err)
		os.Exit(1)
		//log.Fatalf("error al iniciar el servidor %v", err)
	}

}
