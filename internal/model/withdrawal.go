package model

type Withdraw struct {
	OrderNumber string  `json:"order"`
	UserID      int     `json:"user_id"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type ReqWithdraw struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}
