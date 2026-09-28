// Package server creates and start server
package server

import (
	"net/http"
	"time"

	
)

// Server struct dependencies -> Encapsular al servidor http y sus dependencias
type Server struct {
	httpServer *http.Server   // Configuration -> Guarda la configuracion y comportamiento del http server real
	mux        *http.ServeMux // Enrutador -> que handle responde a cada ruta
}

// New create -> para crear un servidor nuevo
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

// Start server -> Metodo para arrancar el servidor http
func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

// Metodo para regitrar un handler en una ruta
func (server *Server) RegisterRoutes(pattern string, handler http.Handler) {
	server.mux.Handle(pattern, handler)
}

// ServeHTTP server
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.mux.ServeHTTP(writer, request)
}
