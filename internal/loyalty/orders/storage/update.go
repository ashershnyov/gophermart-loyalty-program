package storage

import (
	"context"
	"fmt"
)

const qInsertOrder = `
	INSERT INTO gophermart.orders (user_id, number)
	VALUES ($1, $2);
`

// AddOrder inserts a newly created order to the orders table.
// Keeps the `status`, `created` and `updated` according to the table defaults.
func (s *Storage) AddOrder(ctx context.Context, number string, userID int64) error {
	_, err := s.db.ExecContext(ctx, qInsertOrder, userID, number)
	if err != nil {
		return fmt.Errorf("an error occurred when inserting order %s", number)
	}
	return nil
}

const qUpdateOrderStatus = `
	UPDATE gophermart.orders
	SET updated = NOW(),
		status = $1,
		accrual = $2
	WHERE number = $3;
`

// UpdateOrderAccrual updates status and accrual for the order.
func (s *Storage) UpdateOrderAccrual(ctx context.Context, status string, accrual float64, number string) error {
	_, err := s.db.ExecContext(ctx, qUpdateOrderStatus, status, accrual, number)
	if err != nil {
		return fmt.Errorf("an error occurred when updating order %s: %w", number, err)
	}
	return nil
}
