package handlers

import (
	httphelper "Lolopenza/CRUD-F/internal/http-helper"
	"Lolopenza/CRUD-F/internal/models"
	"Lolopenza/CRUD-F/internal/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up!")
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		var user models.User

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			if r.Header.Get("Content-Type") != "application/json" {
				httphelper.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
				return
			}
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			httphelper.WriteError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		id, err := repository.CreateUser(db, user.Email, user.Name, user.Surname)
		if err != nil {
			httphelper.WriteError(w, http.StatusInternalServerError, "Cannot create user")
			log.Println("createUser error:", err)
			return
		}

		err = httphelper.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
		if err != nil {
			log.Println("writeJSON failed:", err)
			return
		}

	}
}

func RecieveAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User

		users, err := repository.GetAllUsers(db)
		if err != nil {
			httphelper.WriteError(w, http.StatusInternalServerError, "error on getting db side")
			return
		}

		err = httphelper.WriteJSON(w, http.StatusOK, users)
		if err != nil {
			log.Println("writeJSON failed:", err)
			return
		}

	}
}

func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		vars := mux.Vars(r)
		id := vars["id"]

		num_id, err := strconv.Atoi(id)
		if err != nil || num_id < 1 {
			httphelper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}

		user, err = repository.GettingUser(db, num_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httphelper.WriteError(w, http.StatusNotFound, "user not found")
				return
			} else {
				httphelper.WriteError(w, http.StatusInternalServerError, "server issue")
				return
			}
		}

		err = httphelper.WriteJSON(w, http.StatusOK, user)
		if err != nil {
			log.Println("writeJSON failed:", err)
			return
		}

	}
}

func ChangeUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		var user models.User

		vars := mux.Vars(r)
		id := vars["id"]

		num_id, err := strconv.Atoi(id)
		if err != nil || num_id < 1 {
			httphelper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}

		user, err = repository.GettingUser(db, num_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httphelper.WriteError(w, http.StatusNotFound, "user not found")
				return
			} else {
				httphelper.WriteError(w, http.StatusInternalServerError, "server issue")
				return
			}
		}

		var input struct {
			Email   string `json:"email"`
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			if r.Header.Get("Content-Type") != "application/json" {
				httphelper.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
				return
			}
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			httphelper.WriteError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		user.Email = input.Email
		user.Name = input.Name
		user.Surname = input.Surname

		user, err = repository.UpdateUser(db, num_id, user.Email, user.Name, user.Surname)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httphelper.WriteError(w, http.StatusNotFound, "user not found")
				return
			} else {
				httphelper.WriteError(w, http.StatusInternalServerError, "server issue")
				return
			}
		}

		err = httphelper.WriteJSON(w, http.StatusOK, user)
		if err != nil {
			log.Println("writeJSON failed:", err)
			return
		}
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		num_id, err := strconv.Atoi(id)
		if err != nil || num_id < 1 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = repository.DeleteUser(db, num_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httphelper.WriteError(w, http.StatusNotFound, "user not found")
				return
			}
			httphelper.WriteError(w, http.StatusInternalServerError, "server issue")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
