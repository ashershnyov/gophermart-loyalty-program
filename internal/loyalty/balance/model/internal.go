package model

import "time"

// IntBalanceInfo is a user's balance info.
type IntBalanceInfo struct {
	Current   float64 `db:"current"`
	Withdrawn float64 `db:"total_withdrawn"`
}

// ToExternal converts from internal data model to external.
func (bi *IntBalanceInfo) ToExternal() BalanceResp {
	return BalanceResp{
		Current:   bi.Current,
		Withdrawn: bi.Withdrawn,
	}
}

// IntWithdrawal represents single withdrawal internally.
type IntWithdrawal struct {
	Created time.Time `db:"created"`
	OrderID string    `db:"order_id"`
	Amount  float64   `db:"amount"`
}

// ToExternal converts from internal data model to external.
func (iw *IntWithdrawal) ToExternal() Withdrawal {
	return Withdrawal{
		ProcessedAt: iw.Created,
		Order:       iw.OrderID,
		Sum:         iw.Amount,
	}
}
