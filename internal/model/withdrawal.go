package model

type Withdraw struct {
	Order       string  `json:"order"`
	UserId      int     `json:"user_id"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type ReqWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
