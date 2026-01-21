package handlers

import (
	httphelper "Lolopenza/CRUD-F/internal/http-helper"
	"Lolopenza/CRUD-F/internal/models"
	"Lolopenza/CRUD-F/internal/service"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	Service *service.UserService
	Logger  *slog.Logger
}

func NewUserHandler(svc *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		Service: svc,
		Logger:  logger,
	}
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up!")
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	defer r.Body.Close()

	var user models.User

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		httphelper.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.Logger.Warn("Invalid JSON", "err", err)
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, err := h.Service.CreateUser(ctx, user.Email, user.Name, user.Surname)
	if err != nil {
		h.Logger.Error("createUser error", "err", err)
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = httphelper.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
	if err != nil {
		h.Logger.Error("writeJSON failed", "err", err)
		return
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	users, err := h.Service.GetAllUsers(ctx)
	if err != nil {
		h.Logger.Error("failed to get users", "err", err)
		httphelper.WriteError(w, http.StatusInternalServerError, "failed to get users")
		return
	}

	if err := httphelper.WriteJSON(w, http.StatusOK, users); err != nil {
		h.Logger.Warn("writeJSON failed", "err", err)
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		h.Logger.Warn("invalid id", "err", err)
		httphelper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.Service.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Warn("user not found", "err", err)
			httphelper.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		h.Logger.Error("failed to get user", "err", err)
		httphelper.WriteError(w, http.StatusInternalServerError, "server issue")
		return
	}

	if err := httphelper.WriteJSON(w, http.StatusOK, user); err != nil {
		h.Logger.Warn("writeJSON failed", "err", err)
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		h.Logger.Warn("invalid id", "err", err)
		httphelper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		httphelper.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Logger.Warn("Invalid JSON", "err", err)
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.Service.UpdateUser(ctx, id, input.Email, input.Name, input.Surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Warn("user not found", "err", err)
			httphelper.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		h.Logger.Error("failed to update user", "err", err)
		httphelper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := httphelper.WriteJSON(w, http.StatusOK, user); err != nil {
		h.Logger.Warn("writeJSON failed", "err", err)
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		h.Logger.Warn("invalid id", "err", err)
		httphelper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Service.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Warn("user not found", "err", err)
			httphelper.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		h.Logger.Error("failed to delete user", "err", err)
		httphelper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
