package model

import "time"

// IntOrder describes an internal single order.
type IntOrder struct {
	UploadedAt time.Time `db:"created"`
	Number     string    `db:"number"`
	Status     string    `db:"status"`
	Accrual    float64   `db:"accrual"`
	UserID     int64     `db:"user_id"`
}

// ToExternal converts from internal data model to external.
func (io *IntOrder) ToExternal() Order {
	return Order{
		UploadedAt: io.UploadedAt,
		Number:     io.Number,
		Status:     io.Status,
		Accrual:    io.Accrual,
		UserID:     io.UserID,
	}
}
