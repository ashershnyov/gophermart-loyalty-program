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

const qGetSingleOrder = `
	SELECT status, accrual, updated, user_id FROM gophermart.orders
	WHERE number = $1;
`

// GetSingleOrder gets a single order by its number.
func (s *Storage) GetSingleOrder(ctx context.Context, number string) (model.IntOrder, error) {
	row := s.db.QueryRowContext(ctx, qGetSingleOrder, number)
	var order = model.IntOrder{}
	err := row.Scan(&order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	if err != nil {
		return model.IntOrder{}, fmt.Errorf("an error occurred when fetching order DB %s: %w", number, err)
	}
	return order, nil
}

const qFindOrdersToPoll = `
	SELECT number, status, accrual, created, user_id FROM gophermart.orders
	WHERE status != $1;
`

// FindOrdersToPoll returns all orders that have to be polled.
func (s *Storage) FindOrdersToPoll(ctx context.Context) ([]model.IntOrder, error) {
	rows, err := s.db.QueryContext(ctx, qFindOrdersToPoll, model.StatusProcessed)
	if err != nil {
		return nil, fmt.Errorf(errFetchingOrders, err)
	}
	defer rows.Close()

	orders := []model.IntOrder{}
	for rows.Next() {
		var order model.IntOrder
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
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
