package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/service"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/middleware"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type ordersService interface {
	GetOrders(ctx context.Context, userID int64) (model.GetOrdersResp, error)
	AddOrder(ctx context.Context, number string, userID int64) error
	GetSingleOrder(ctx context.Context, number string, userID int64) (*model.Order, error)
}

// Handler is a orders handler.
type Handler struct {
	os      ordersService
	withJwt middleware.Middleware
}

// New creates a new order handler.
func New(service ordersService, withJwt middleware.Middleware) *Handler {
	return &Handler{
		os:      service,
		withJwt: withJwt,
	}
}

func (h *Handler) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(jwt.CtxKey).(int64)
		if !ok {
			http.Error(w, "no user id in context", http.StatusInternalServerError)
			return
		}

		orders, err := h.os.GetOrders(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(orders) < 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func (h *Handler) addOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(jwt.CtxKey).(int64)
		if !ok {
			http.Error(w, "no user id in context", http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		orderNumber := string(buf.Bytes())

		_, err = h.os.GetSingleOrder(r.Context(), orderNumber, userID)
		if errors.Is(err, service.ErrInvalidOrderNumber) {
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, service.ErrOrderNumberTaken) {
			http.Error(w, "order with that number already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrOrderAlreadyExists) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if err != nil && !errors.Is(err, service.ErrNoOrderFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.os.AddOrder(r.Context(), orderNumber, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// RegisterRoutes registers routes of orders logic.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/api/user/orders", h.withJwt(h.getOrders()).ServeHTTP)
	router.Post("/api/user/orders", h.withJwt(h.addOrder()).ServeHTTP)
}
