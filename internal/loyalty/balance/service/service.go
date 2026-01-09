package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/storage"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/luhn"
)

var (
	// ErrInvalidOrderNumber means the order number is not valid.
	ErrInvalidOrderNumber = errors.New("invalid order number")
	// ErrNotEnoughPoints means that user doesn't have enough points to perform a withdrawal.
	ErrNotEnoughPoints = errors.New("not enough points")
)

type balanceStorage interface {
	GetBalance(ctx context.Context, userID int64) (model.IntBalanceInfo, error)
	GetWithdrawals(ctx context.Context, userID int64) ([]model.IntWithdrawal, error)
	Withdraw(ctx context.Context, userID int64, orderID string, amount float64) error
	AddAccrual(ctx context.Context, userID int64, amount float64) error
}

// BalanceService is the service layer for balance logic.
type BalanceService struct {
	storage balanceStorage
}

// New creates a new balance service.
func New(db db.DB) BalanceService {
	return BalanceService{
		storage: storage.NewStorage(db),
	}
}

// GetBalance returns balance information for user with the provided userID.
func (bs *BalanceService) GetBalance(ctx context.Context, userID int64) (model.BalanceResp, error) {
	balance, err := bs.storage.GetBalance(ctx, userID)
	if err != nil {
		return model.BalanceResp{}, fmt.Errorf("error getting balance: %w", err)
	}
	return balance.ToExternal(), nil
}

// GetWithdrawals withdrawals information for user with the provided userID.
func (bs *BalanceService) GetWithdrawals(ctx context.Context, userID int64) (model.GetWithdrawalsResp, error) {
	withdrawals, err := bs.storage.GetWithdrawals(ctx, userID)
	if err != nil {
		return make(model.GetWithdrawalsResp, 0), fmt.Errorf("error getting withdrawals: %w", err)
	}

	resp := make(model.GetWithdrawalsResp, len(withdrawals), 0)
	for _, w := range withdrawals {
		resp = append(resp, w.ToExternal())
	}

	return resp, nil
}

// Withdraw performs a withdrawal operation.
func (bs *BalanceService) Withdraw(ctx context.Context, withdrawal model.WithdrawReq, userID int64) error {
	if !luhn.Validate(withdrawal.Order) {
		return ErrInvalidOrderNumber
	}

	balance, err := bs.storage.GetBalance(ctx, userID)
	if err != nil {
		slog.Warn(err.Error())
		return fmt.Errorf("error getting balance: %w", err)
	}
	if withdrawal.Sum > balance.Current {
		return ErrNotEnoughPoints
	}

	err = bs.storage.Withdraw(ctx, userID, withdrawal.Order, withdrawal.Sum)
	if err != nil {
		slog.Warn(err.Error())
		return fmt.Errorf("error performig withdrawal: %w", err)
	}
	return nil
}

// AddAccrual adds passed amount to user's balance.
func (bs *BalanceService) AddAccrual(ctx context.Context, userID int64, amount float64) error {
	return bs.storage.AddAccrual(ctx, userID, amount)
}
