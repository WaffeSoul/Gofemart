package model

type Order struct {
	Number     string `json:"number"`
	UserId     int    `json:"user_id"`
	UploadedAt string `json:"uploaded_at"`
}

type OrderWithAccrual struct {
	Status     string  `json:"status"`
	Number     string  `json:"number"`
	UploadedAt string  `json:"uploaded_at"`
	Accrual    float64 `json:"accrual"`
}

func OrderToOrderWithAccrual(order Order, accrual float64, status string) OrderWithAccrual {
	switch status {
	case "REGISTERED":
		status = "NEW"
	case "INVALID":
		status = "INVALID"
	case "PROCESSING":
		status = "PROCESSING"
	case "PROCESSED":
		status = "PROCESSED"
	}
	return OrderWithAccrual{
		Status:     status,
		Number:     order.Number,
		UploadedAt: order.UploadedAt,
		Accrual:    accrual,
	}
}

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
