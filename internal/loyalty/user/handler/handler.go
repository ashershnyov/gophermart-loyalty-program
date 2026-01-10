package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/service"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	Login(ctx context.Context, req model.LoginReq) (string, error)
	Register(ctx context.Context, login, password string) (string, error)
}

// Handler is a user state handler.
type Handler struct {
	us userService
}

// New creates a new user handler.
func New(service userService) *Handler {
	return &Handler{
		us: service,
	}
}

func (h *Handler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userData model.RegisterReq
		if err := json.Unmarshal(buf.Bytes(), &userData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tok, err := h.us.Register(r.Context(), userData.Login, userData.Password)
		if err != nil {
			if errors.Is(err, service.ErrDuplicateLogin) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", "Bearer "+tok)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userData model.LoginReq
		if err := json.Unmarshal(buf.Bytes(), &userData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tok, err := h.us.Login(r.Context(), userData)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrWrongPassword) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", "Bearer "+tok)
		w.WriteHeader(http.StatusOK)
	}
}

// RegisterRoutes registers routes of user logic.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/api/user/register", h.register())
	router.Post("/api/user/login", h.login())
}
