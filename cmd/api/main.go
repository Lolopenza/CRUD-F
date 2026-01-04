package main

import (
	"Lolopenza/CRUD-F/internal/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	dsn = "postgres://anvar:1234@localhost:5432/crud_api?sslmode=disable"
)

func main() {
	router := mux.NewRouter()

	db := DB_init()

	router.HandleFunc("/healthcheck", handlers.HealthcheckHandler).Methods("GET")

	router.HandleFunc("/users", handlers.CreateUserHandler(db)).Methods("POST")
	router.HandleFunc("/users", handlers.RecieveAllUsersHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.GetUserHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.ChangeUserHandler(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", (handlers.DeleteUserHandler(db))).Methods("DELETE")

	defer db.Close()

	fmt.Println("Server starting on port :3838")
	log.Fatal(http.ListenAndServe(":3838", router))
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
