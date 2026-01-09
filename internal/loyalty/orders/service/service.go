package service

import (
	"context"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/storage"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
)

type ordersStorage interface {
	GetOrders(ctx context.Context, userID int64) ([]model.IntOrder, error)
	AddOrder(ctx context.Context, order model.IntOrder, userID int64) error
}

// OrdersService is the service layer for orders logic.
type OrdersService struct {
	storage ordersStorage
}

// New creates a new orders service.
func New(db db.DB) OrdersService {
	return OrdersService{
		storage: storage.NewStorage(db),
	}
}

// GetOrders returns orders of a user with the specified userID.
func (os *OrdersService) GetOrders(ctx context.Context, userID int64) (model.GetOrdersResp, error) {
	orders, err := os.storage.GetOrders(ctx, userID)
	if err != nil {
		return make(model.GetOrdersResp, 0), fmt.Errorf("error while getting orders: %w", err)
	}

	resp := make(model.GetOrdersResp, len(orders), 0)
	for _, o := range orders {
		resp = append(resp, o.ToExternal())
	}

	return resp, nil
}

// AddOrder adds order for a user with the specified userID.
func (os *OrdersService) AddOrder(ctx context.Context, order *model.Order, userID int64) error {
	err := os.storage.AddOrder(ctx, order.ToInternal(), userID)
	if err != nil {
		return fmt.Errorf("error while adding order: %w", err)
	}
	return nil
}
