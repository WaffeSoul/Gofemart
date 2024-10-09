package model

type Withdraw struct {
	Order       string `json:"order"`
	Sum         int    `json:"sum"`
	UserId      int    `json:"user_id"`
	ProcessedAt string `json:"processed_at"`
}

type ReqWithdraw struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}
