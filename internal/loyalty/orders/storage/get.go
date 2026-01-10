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
	var orders = []model.IntOrder{}
	err := s.db.QueryManyContext(ctx, &orders, qGetOrders, userID)
	if err != nil {
		return nil, fmt.Errorf(errFetchingOrders, err)
	}
	return orders, nil
}

const qGetSingleOrder = `
	SELECT status, accrual, user_id FROM gophermart.orders
	WHERE number = $1;
`

// GetSingleOrder gets a single order by its number.
func (s *Storage) GetSingleOrder(ctx context.Context, number string) (model.IntOrder, error) {
	var order = model.IntOrder{}
	err := s.db.QueryOneContext(ctx, &order, qGetSingleOrder, number)
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
	orders := []model.IntOrder{}
	err := s.db.QueryManyContext(ctx, &orders, qFindOrdersToPoll, model.StatusProcessed)
	if err != nil {
		return nil, fmt.Errorf(errFetchingOrders, err)
	}
	return orders, nil

}
