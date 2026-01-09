package handler

import (
	"context"
	"net/http"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type ordersService interface {
	GetOrders(ctx context.Context, userID int64) (model.GetOrdersResp, error)
	AddOrder(ctx context.Context, order *model.Order, userID int64) error
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
	panic("Not implmeneted")
}

func (h *Handler) addOrders() http.HandlerFunc {
	panic("Not implmeneted")
}

// RegisterRoutes registers routes of orders logic.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/api/user/orders ", h.withJwt(h.addOrders()).ServeHTTP)
	router.Get("/api/user/orders", h.withJwt(h.getOrders()).ServeHTTP)
}
