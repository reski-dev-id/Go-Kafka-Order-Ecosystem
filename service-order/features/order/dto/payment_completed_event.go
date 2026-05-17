package dto

type PaymentCompletedEvent struct {
	OrderID string  `json:"orderId"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}
