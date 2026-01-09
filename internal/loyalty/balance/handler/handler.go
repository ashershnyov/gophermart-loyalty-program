package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/service"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/middleware"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type balanceService interface {
	GetBalance(ctx context.Context, userID int64) (model.BalanceResp, error)
	Withdraw(ctx context.Context, withdrawal model.WithdrawReq, userID int64) error
	GetWithdrawals(ctx context.Context, userID int64) (model.GetWithdrawalsResp, error)
}

// Handler is a user balance handler.
type Handler struct {
	bs      balanceService
	withJwt middleware.Middleware
}

// New returns a newly created user balance handler.
func New(service balanceService, withJwt middleware.Middleware) *Handler {
	return &Handler{
		bs:      service,
		withJwt: withJwt,
	}
}

func (h *Handler) getBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(jwt.CtxKey).(int64)
		if !ok {
			http.Error(w, "no user id in context", http.StatusInternalServerError)
			return
		}

		bal, err := h.bs.GetBalance(r.Context(), userID)
		if err != nil {
			http.Error(w, "error getting balance", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(bal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func (h *Handler) withdraw() http.HandlerFunc {
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

		var userData model.WithdrawReq
		if err := json.Unmarshal(buf.Bytes(), &userData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.bs.Withdraw(r.Context(), userData, userID)
		if err != nil {
			if errors.Is(err, service.ErrNotEnoughPoints) {
				http.Error(w, "not enough loyalty points", http.StatusPaymentRequired)
				return
			}
			if errors.Is(err, service.ErrInvalidOrderNumber) {
				http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Accept", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) getWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(jwt.CtxKey).(int64)
		if !ok {
			http.Error(w, "no user id in context", http.StatusInternalServerError)
			return
		}

		withdrawals, err := h.bs.GetWithdrawals(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(withdrawals) < 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

// RegisterRoutes registers routes of balance logic.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/api/user/balance", h.withJwt(h.getBalance()).ServeHTTP)
	router.Post("/api/user/balance/withdraw", h.withJwt(h.withdraw()).ServeHTTP)
	router.Get("/api/user/withdrawals", h.withJwt(h.getWithdrawals()).ServeHTTP)
}
