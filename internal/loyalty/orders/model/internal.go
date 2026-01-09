package model

import "time"

// IntOrder describes an internal single order.
type IntOrder struct {
	UploadedAt time.Time
	Number     string
	Status     string
	Accrual    float64
	UserID     int64
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
