package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/storage"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/luhn"
)

var (
	// ErrInvalidOrderNumber means the order number is not valid.
	ErrInvalidOrderNumber = errors.New("invalid order number")
	// ErrOrderAlreadyExists means the order with a certain number already exists for the specified userID.
	ErrOrderAlreadyExists = errors.New("order already exists")
	// ErrOrderNumberTaken means the order number is already taken.
	ErrOrderNumberTaken = errors.New("order number taken")
	// ErrNoOrderFound means the order was not found for the specfic number.
	ErrNoOrderFound = errors.New("no order found")
)

type ordersStorage interface {
	GetOrders(ctx context.Context, userID int64) ([]model.IntOrder, error)
	AddOrder(ctx context.Context, number string, userID int64) error
	GetSingleOrder(ctx context.Context, number string) (model.IntOrder, error)
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

	resp := make(model.GetOrdersResp, 0, len(orders))
	for _, o := range orders {
		resp = append(resp, o.ToExternal())
	}

	return resp, nil
}

// AddOrder adds order for a user with the specified userID.
func (os *OrdersService) AddOrder(ctx context.Context, number string, userID int64) error {
	if !luhn.Validate(number) {
		return ErrInvalidOrderNumber
	}
	err := os.storage.AddOrder(ctx, number, userID)
	if err != nil {
		return fmt.Errorf("error while adding order: %w", err)
	}
	return nil
}

// GetSingleOrder gets a single order by its number and userID.
// If an order exists but userID doesn't match will return ErrOrderAlreadyExists.
func (os *OrdersService) GetSingleOrder(ctx context.Context, number string, userID int64) (*model.Order, error) {
	order, err := os.storage.GetSingleOrder(ctx, number)
	if userID == order.UserID && err == nil {
		return nil, ErrOrderAlreadyExists
	}
	if userID != order.UserID && err == nil {
		return nil, ErrOrderNumberTaken
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoOrderFound
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching order: %w", err)
	}
	extOrder := order.ToExternal()
	return &extOrder, nil
}
