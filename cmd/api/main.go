package main

import (
	"Lolopenza/CRUD-F/internal/handlers"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	dsn = "postgres://anvar:1234@localhost:5432/crud_api?sslmode=disable"
)

func main() {
	db := DB_init()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", handlers.HealthcheckHandler).Methods("GET")

	router.HandleFunc("/users", handlers.CreateUserHandler(db)).Methods("POST")
	router.HandleFunc("/users", handlers.RecieveAllUsersHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.GetUserHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.ChangeUserHandler(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", (handlers.DeleteUserHandler(db))).Methods("DELETE")

	srv := &http.Server{
		Addr:         ":3838",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server started on :3838")
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

	if err := db.Close(); err != nil {
		log.Println("db close error:", err)
	}

	log.Println("server exited gracefully")
}

func DB_init() *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to db successfully!")
	}

	return db
}

//"postgres://anvar:1234@localhost:5432/crud_api"
