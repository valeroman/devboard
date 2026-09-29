// Package main api
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/valeroman/devboard/internal/handler"
	"github.com/valeroman/devboard/internal/middlewares"
	"github.com/valeroman/devboard/internal/server"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Aqui arranca el servidor
	srv := server.New(":8090")

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	healthHander := handler.NewHealtHandler()

	srv.RegisterRoutes("GET /health", healthHander)
	srv.RegisterRoutes("GET /ready", healthHander)

	// Ruta temporal para probar Recovery
	srv.RegisterRoutes("GET /panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error provocado para probar Recovery")
	}))

	//log.Println("Servidor iniciado en :8090")
	logger.Info("servidor iniciado", slog.String("addr", ":8090"))

	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}

}
