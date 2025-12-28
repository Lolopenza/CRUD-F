package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	dsn = "postgres://anvar:1234@localhost:5432/crud_api?sslmode=disable"
)

type User struct {
	User_ID    int       `json:"usr_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up!")
}

func createuserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		id, err := createUser(db, user.Email, user.Name, user.Surname)
		if err != nil {
			http.Error(w, "Cannot create user", http.StatusInternalServerError)
			log.Println("createUser error:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User created successfully with id %d\n", id)
	}

}

func createUser(db *sql.DB, email, name, surname string) (int, error) {
	var id int

	stmt := `INSERT INTO users ( email, name, surname) VALUES ($1, $2, $3) RETURNING usr_id`

	err := db.QueryRow(stmt, email, name, surname).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func recieveallusersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "recieveallusersHandler")
}

func getuserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getuserHandler")
}

func changeuserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "changeuserHandler")
}

func deleteuserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "deleteuserHandler")
}

func main() {
	router := mux.NewRouter()

	db := DB_init()

	router.HandleFunc("/healthcheck", healthcheckHandler).Methods("GET")

	router.HandleFunc("/users", createuserHandler(db)).Methods("POST")
	router.HandleFunc("/users", recieveallusersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", getuserHandler).Methods("GET")
	router.HandleFunc("/users/{id}", changeuserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteuserHandler).Methods("DELETE")

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
