package main

import (
	"Lolopenza/CRUD-F/internal/config"
	"Lolopenza/CRUD-F/internal/db"
	"Lolopenza/CRUD-F/internal/handlers"
	"Lolopenza/CRUD-F/internal/repository"
	"Lolopenza/CRUD-F/internal/service"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	).With(
		slog.String("service", "crud-api"),
	)

	slog.SetDefault(logger)

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	logger.Info("Config loaded")

	// DB
	database := db.New(cfg.DB.DSN, logger)
	defer database.Close()

	// DI: repository → service → handler
	repo := repository.New(database)
	service := service.NewUserService(repo)
	handler := handlers.NewUserHandler(service, logger)

	// Router
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", handlers.HealthcheckHandler).Methods("GET")

	router.HandleFunc("/users", handler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info(
			"server started",
			"port", srv.Addr,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(
				"listen error",
				"err", err,
			)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "err", err)
	}

	logger.Info("server exited gracefully")
}
