// Package main api
package main

import (
	"log"

	"github.com/valeroman/devboard/internal/handler"
	"github.com/valeroman/devboard/internal/server"
)

func main() {
	// Aqui arranca el servidor
	srv := server.New(":8090")

	healthHander := handler.NewHealtHandler()

	srv.RegisterRoutes("GET /health", healthHander)
	srv.RegisterRoutes("GET /ready", healthHander)

	log.Println("Servidor iniciado en :8090")

	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}

}
