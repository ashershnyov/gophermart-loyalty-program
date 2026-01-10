package storage

import (
	"context"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/model"
)

const errFetchingWithdrawals = "an error occurred when fetching withdrawals from DB: %w"

const qGetBalance = `
	SELECT current, total_withdrawn FROM gophermart.balance
	WHERE user_id=$1;
`

// GetBalance returns balance for user with the provided userID.
func (s *Storage) GetBalance(ctx context.Context, userID int64) (model.IntBalanceInfo, error) {
	var balance = model.IntBalanceInfo{}
	err := s.db.QueryOneContext(ctx, balance, qGetBalance, userID)
	if err != nil {
		return model.IntBalanceInfo{}, fmt.Errorf("an error occurred when fetching balance DB for user %v: %w", userID, err)
	}
	return balance, nil
}

const qGetWithdrawals = `
	SELECT order_id, amount, created FROM gophermart.withdrawals
	WHERE user_id=$1
	ORDER BY created DESC;
`

// GetWithdrawals returns withdrawals for user with the provided userID.
func (s *Storage) GetWithdrawals(ctx context.Context, userID int64) ([]model.IntWithdrawal, error) {
	var withdrawals = []model.IntWithdrawal{}
	err := s.db.QueryManyContext(ctx, withdrawals, qGetWithdrawals, userID)
	if err != nil {
		return nil, fmt.Errorf(errFetchingWithdrawals, err)
	}
	return withdrawals, nil
}
