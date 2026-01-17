package main

import (
	"Lolopenza/CRUD-F/internal/config"
	"Lolopenza/CRUD-F/internal/db"
	"Lolopenza/CRUD-F/internal/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.Load()

	database := db.New(cfg.DB.DSN)
	defer database.Close()

	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", handlers.HealthcheckHandler).Methods("GET")

	router.HandleFunc("/users", handlers.CreateUserHandler(database)).Methods("POST")
	router.HandleFunc("/users", handlers.RecieveAllUsersHandler(database)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.GetUserHandler(database)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.ChangeUserHandler(database)).Methods("PUT")
	router.HandleFunc("/users/{id}", (handlers.DeleteUserHandler(database))).Methods("DELETE")

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server started on :" + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	if err := database.Close(); err != nil {
		log.Println("db close error:", err)
	}

	log.Println("server exited gracefully")
}
