// Package server creates and start server
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

// Server struct dependencies -> Encapsular al servidor http y sus dependencies
type Server struct {
	httpServer *http.Server   // Configuration -> Guarda la configuration y comportamiento del http server real
	mux        *http.ServeMux // Enrutador -> que handle responde a cada ruta
	logger     *slog.Logger
}

// New create -> para crear un servidor nuevo
func New(addr string, opts ...Option) *Server {

	cfg := defaultConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	mux := http.NewServeMux() // Crea un nuevo multiplexor o enrutador http

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       cfg.readTimeout,
		WriteTimeout:      cfg.writeTimeout,
		IdleTimeout:       cfg.idleTimeout,
		ReadHeaderTimeout: cfg.readHeaderTimeout,
	} // Crear una instancia del servidor http

	return &Server{
		httpServer: httpServer,
		mux:        mux,
		logger:     cfg.logger,
	}
}

// Start server -> Metodo para arrancar el servidor http
func (server *Server) Start() error {

	server.logger.Info("iniciado servidor...", slog.String("addr", server.httpServer.Addr))

	err := server.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server http: %w", err)
	}

	return nil
}

// Shutdown method
func (server *Server) Shutdown(ctx context.Context) error {

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("error en shutdown", slog.Any("error", err))
		return err
	}

	server.logger.Info("Servidor detenido correctamente")

	return nil
}

// RegisterRoutes Metodo para regitrar un handler en una ruta
func (server *Server) RegisterRoutes(pattern string, handler http.Handler) {
	server.mux.Handle(pattern, handler)
}

// ServeHTTP server
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.mux.ServeHTTP(writer, request)
}

// Use Metodo
func (server *Server) Use(middleware func(http.Handler) http.Handler) {
	server.httpServer.Handler = middleware(server.httpServer.Handler)
}

// UseChain cadena para middleware
// func (server *Server) UseChain(middlewaresParams ...middlewares.Middleware) {
// 	server.httpServer.Handler = middlewares.Chain(server.httpServer.Handler, middlewaresParams...)
// }
