package model

import "time"

// Order describes a single order.
type Order struct {
	UploadedAt time.Time `json:"uploaded_at"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UserID     int64     `json:"-"`
}

// ToExternal converts from external data model to internal.
func (o *Order) ToInternal() IntOrder {
	return IntOrder{
		UploadedAt: o.UploadedAt,
		Number:     o.Number,
		Status:     o.Status,
		Accrual:    o.Accrual,
	}
}

// GetOrdersResp is a response body for getOrders handler.
type GetOrdersResp []Order

// AccrualOrder defines an accrual response body.
type AccrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
	UserID  int64   `json:"-"`
}
