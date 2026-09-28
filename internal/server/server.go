package server

import (
	"net/http"
	"time"
)

// Encapsular al servidor http y sus dependencias
type Server struct {
	httpServer *http.Server   // Guarda la configuracion y comportamiento del http server real
	mux        *http.ServeMux // Enrutador -> que handle responde a cada ruta
}

// funcion para crear un servidor nuevo
func New(addr string) *Server {
	mux := http.NewServeMux() // Crea un nuevo multiplexor o enrutador http

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	} // Crear una instancia del servidor http

	return &Server{
		httpServer: httpServer,
		mux:        mux,
	}
}

// Metodo para arrancar el servidor http
func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}
