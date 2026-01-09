package storage

import (
	"context"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
)

const errFetchingOrders = "an error occurred when fetching orders from DB: %w"

const qGetOrders = `
	SELECT number, status, accrual, created
	FROM gophermart.orders
	WHERE user_id=$1;
`

// GetOrders returns all the orders for the provided userID.
func (s *Storage) GetOrders(ctx context.Context, userID int64) ([]model.IntOrder, error) {
	rows, err := s.db.QueryContext(ctx, qGetOrders, userID)
	if err != nil {
		return nil, fmt.Errorf(errFetchingOrders, err)
	}
	defer rows.Close()

	orders := []model.IntOrder{}
	for rows.Next() {
		var order model.IntOrder
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf(errFetchingOrders, err)
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(errFetchingOrders, err)
	}

	return orders, nil
}
