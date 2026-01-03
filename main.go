package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	Created_At time.Time `json:"created_at,omitempty"`
	Updated_At time.Time `json:"updated_at,omitempty"`
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

func recieveallusersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var users []User

		users, err := getAllUsers(db)
		if err != nil {
			http.Error(w, "error on getting db side", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	}
}

func getAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT usr_id, email, name, surname FROM users`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.User_ID, &u.Email, &u.Name, &u.Surname); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func getuserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var user User

		vars := mux.Vars(r)
		id := vars["id"]

		num_id, err := strconv.Atoi(id)
		if err != nil || num_id < 1 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		user, err = gettingUser(db, num_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, " server issue", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

	}
}

func gettingUser(db *sql.DB, num_id int) (User, error) {
	var u User

	query := `SELECT usr_id, email, name, surname
			FROM users
			WHERE usr_id = $1`

	err := db.QueryRow(query, num_id).Scan(&u.User_ID, &u.Email, &u.Name, &u.Surname)
	if err != nil {
		return u, err
	}

	return u, nil
}

func changeuserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var user User

		vars := mux.Vars(r)
		id := vars["id"]

		num_id, err := strconv.Atoi(id)
		if err != nil || num_id < 1 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		user, err = gettingUser(db, num_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, " server issue", http.StatusInternalServerError)
				return
			}
		}

		var input struct {
			Email   string `json:"email"`
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		user.Email = input.Email
		user.Name = input.Name
		user.Surname = input.Surname

		user, err = updateUser(db, num_id, user.Email, user.Name, user.Surname)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, " server issue", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	}
}

func updateUser(db *sql.DB, num_id int, email, name, surname string) (User, error) {
	var u User

	query := `
		UPDATE users
		SET email = $1, name = $2, surname = $3, updated_at = now()
		WHERE usr_id = $4
		RETURNING usr_id, email, name, surname, created_at, updated_at
	`

	err := db.QueryRow(query, email, name, surname, num_id).
		Scan(&u.User_ID, &u.Email, &u.Name, &u.Surname, &u.Created_At, &u.Updated_At)

	if err != nil {
		return u, err
	}

	return u, nil
}

func deleteuserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "deleteuserHandler")
}

func main() {
	router := mux.NewRouter()

	db := DB_init()

	router.HandleFunc("/healthcheck", healthcheckHandler).Methods("GET")

	router.HandleFunc("/users", createuserHandler(db)).Methods("POST")
	router.HandleFunc("/users", recieveallusersHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", getuserHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", changeuserHandler(db)).Methods("PUT")
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
