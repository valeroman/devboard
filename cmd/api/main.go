// Package main api
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/valeroman/devboard/docs"
	"github.com/valeroman/devboard/internal/auth"
	"github.com/valeroman/devboard/internal/handler"
	"github.com/valeroman/devboard/internal/logger"
	"github.com/valeroman/devboard/internal/middlewares"
	pgrepo "github.com/valeroman/devboard/internal/repository/postgres"
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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escribe el token con el formato: Bearer {token}

func main() {

	// 1. Configuración del entorno
	// 2. Logger
	log := setupLogger()

	dbURL := os.Getenv("DATABASE_URL")

	// Crear el pool de conexiones
	pool, err := pgrepo.NewPool(context.Background(), pgrepo.Config{
		URL:             dbURL,
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifeTime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	})

	if err != nil {
		log.Error("no se pudo conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	log.Info("conexion a base de datos establecida")

	// Auth Infrastructure
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshTokenuration := 7 * 24 * time.Hour

	jwtManager, err := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:            []byte(jwtSecret),
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: refreshTokenuration,
	})

	if err != nil {
		log.Error("error configurando JWT", slog.Any("error", err))
		os.Exit(1)
	}

	// 3. Dependencies compartidas

	// Infrastructure: adapters
	// userRepo := memory.NewUserRepository()
	// taskRepo := memory.NewTaskRepository()
	taskRepo := pgrepo.NewTaskRepository(pool)
	workspaceRepo := pgrepo.NewWorkspaceRepository(pool)
	projectRepo := pgrepo.NewProjectRepository(pool)
	commentRepo := pgrepo.NewCommentRepository(pool)

	// Use cases: lógica del negocio
	taskUC := usecase.NewTaskUseCase(taskRepo, taskRepo)
	workspaceUC := usecase.NewWorkspaceUseCase(workspaceRepo)
	projectUC := usecase.NewProjectUseCase(projectRepo, workspaceRepo)
	commentUC := usecase.NewCommentUseCase(commentRepo, taskRepo)

	// Validación compartida
	validate := validator.New()

	// 4. Handlers
	healthHandler := handler.NewHealtHandler()
	taskHandler := handler.NewTaskHandler(taskUC, validate, log)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceUC, validate, log)
	projectHandler := handler.NewProjectHandler(projectUC, validate, log)
	commentHandler := handler.NewCommentHandler(commentUC, validate, log)

	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	// 5. Servidor con Functional options
	// Aqui arranca el servidor
	srv := server.New(":8090",
		server.Withlogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	// Middleware Auth
	authMW := middlewares.RequireAuth(jwtManager, log)

	// 6. Registro de rutas publicas
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

	// 6.1 Registro de rutas protegidas - Tasks
	srv.RegisterRoutes("POST /api/v1/tasks", middlewares.Chain(http.HandlerFunc(taskHandler.Create), authMW))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}", middlewares.Chain(http.HandlerFunc(taskHandler.Get), authMW))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/status", middlewares.Chain(http.HandlerFunc(taskHandler.UpdateStatus), authMW))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/assign", middlewares.Chain(http.HandlerFunc(taskHandler.Assign), authMW))
	srv.RegisterRoutes("GET /api/v1/projects/{id}/tasks", middlewares.Chain(http.HandlerFunc(taskHandler.ListByProject), authMW))

	// 6.2 Registro de rutas protegidas - Workspaces
	srv.RegisterRoutes("POST /api/v1/workspaces", middlewares.Chain(http.HandlerFunc(workspaceHandler.Create), authMW))
	srv.RegisterRoutes("GET /api/v1/workspaces", middlewares.Chain(http.HandlerFunc(workspaceHandler.List), authMW))
	srv.RegisterRoutes("GET /api/v1/workspaces/{id}", middlewares.Chain(http.HandlerFunc(workspaceHandler.Get), authMW))
	srv.RegisterRoutes("POST /api/v1/workspaces/{id}/members", middlewares.Chain(http.HandlerFunc(workspaceHandler.AddMember), authMW))
	srv.RegisterRoutes("GET /api/v1/workspaces/{id}/members", middlewares.Chain(http.HandlerFunc(workspaceHandler.ListMembers), authMW))
	srv.RegisterRoutes("DELETE /api/v1/workspaces/{id}/members/{userID}", middlewares.Chain(http.HandlerFunc(workspaceHandler.RemoveMember), authMW))

	// 6.3 Registro de rutas protegidas - Projects
	srv.RegisterRoutes("POST /api/v1/projects", middlewares.Chain(http.HandlerFunc(projectHandler.Create), authMW))
	srv.RegisterRoutes("GET /api/v1/projects/{id}", middlewares.Chain(http.HandlerFunc(projectHandler.Get), authMW))
	srv.RegisterRoutes("GET /api/v1/workspaces/{id}/projects", middlewares.Chain(http.HandlerFunc(projectHandler.ListByWorkspace), authMW))

	// 6.4 Registro de rutas protegidas - Comments
	srv.RegisterRoutes("POST /api/v1/tasks/{id}/comments", middlewares.Chain(http.HandlerFunc(commentHandler.Create), authMW))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}/comments", middlewares.Chain(http.HandlerFunc(commentHandler.ListByTask), authMW))
	srv.RegisterRoutes("DELETE /api/v1/comments/{id}", middlewares.Chain(http.HandlerFunc(commentHandler.Delete), authMW))

	// 7. Middleware chain
	srv.UseChain(
		middlewares.Recovery(log),
		middlewares.Logger(log),
		middlewares.SecurityHeaders(middlewares.SecurityConfig{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:5173",
			},
		}),
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
