package main

import (
	"log"

	"github.com/valeroman/devboard/internal/server"
)

func main() {
	// Aqui arranca el servidor
	srv := server.New(":8090")

	log.Println("Servidor iniciado en :8090")

	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}
}
