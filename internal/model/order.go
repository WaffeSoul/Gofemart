package model

type Order struct {
	Number     string  `json:"number"`
	UserID     int     `json:"user_id"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func (o *Order) AddAccrual(accrual float64, status string) {
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
	o.Status = status
	o.Accrual = accrual
}

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
