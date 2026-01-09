package storage

import (
	"context"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
)

const qInsertOrder = `
	INSERT INTO gophermart.orders (user_id, number)
	VALUES ($1, $2);
`

// AddOrder inserts a newly created order to the orders table.
// Keeps the `status`, `created` and `updated` according to the table defaults.
func (s *Storage) AddOrder(ctx context.Context, order model.IntOrder, userID int64) error {
	_, err := s.db.ExecContext(ctx, qInsertOrder, userID, order.Number)
	if err != nil {
		return fmt.Errorf("an error occurred when inserting order %s", order.Number)
	}
	return nil
}
