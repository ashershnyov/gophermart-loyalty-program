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
