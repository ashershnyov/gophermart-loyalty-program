package model

import "time"

// IntBalanceInfo is a user's balance info.
type IntBalanceInfo struct {
	Current   float64
	Withdrawn float64
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
	Created time.Time
	OrderID string
	Amount  float64
}

// ToExternal converts from internal data model to external.
func (iw *IntWithdrawal) ToExternal() Withdrawal {
	return Withdrawal{
		ProcessedAt: iw.Created,
		Order:       iw.OrderID,
		Sum:         iw.Amount,
	}
}
