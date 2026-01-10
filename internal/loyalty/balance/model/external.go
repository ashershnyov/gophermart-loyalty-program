package model

import "time"

// BalanceResp is a response body for balance handler.
type BalanceResp struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// WithdrawReq is a request body for withdraw handler.
type WithdrawReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// Withdrawal defines a single withdrawal info.
type Withdrawal struct {
	ProcessedAt time.Time `json:"processed_at"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
}

// GetWithdrawalsResp is a response body for getWithdrawals handler.
type GetWithdrawalsResp []Withdrawal
