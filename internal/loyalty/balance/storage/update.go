package storage

import (
	"context"
	"fmt"
)

const errWithdrawing = "failed to withdraw from user's balance: %w"

const (
	qWithdrawFromBalance = `
		UPDATE gophermart.withdrawals
		SET updated = NOW(),
			current = current + $1,
			total_withdrawn = total_withdrawn - $1
		WHERE user_id = $2;
	`
	qAddWithdrawal = `
		INSER INTO gophermart.withdrawals (user_id, order_id, amount)
		VALUES ($1, $2, $3);
	`
)

// Withdraw adds a new withdrawal and deducts the amount
// from the balance in one transaction for user with the given userID.
func (s *Storage) Withdraw(ctx context.Context, userID int64, orderID string, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(errWithdrawing, err)
	}

	_, err = tx.ExecContext(ctx, qWithdrawFromBalance, amount, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf(errWithdrawing, err)
	}

	_, err = tx.ExecContext(ctx, qAddWithdrawal, userID, orderID, amount)
	if err != nil {
		return fmt.Errorf(errWithdrawing, err)
	}

	return tx.Commit()
}

const qAddAccrual = `
	UPDATE gophermart.balance
	SET updated = NOW(),
		current = current + $1
	WHERE user_id = $2;
`

// AddAccrual adds passed amount to user's balance.
func (s *Storage) AddAccrual(ctx context.Context, userID int64, amount float64) error {
	_, err := s.db.ExecContext(ctx, qAddAccrual, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to add accrual to user's balance: %w", err)
	}
	return nil
}

const qAddBalance = `
	INSERT INTO gophermart.balance (user_id)
	VALUES ($1);
`

// AddBalance adds a new balance row for the provicded userID.
func (s *Storage) AddBalance(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, qAddBalance, userID)
	if err != nil {
		return fmt.Errorf("failed adding balance for new user: %w", err)
	}
	return nil
}
