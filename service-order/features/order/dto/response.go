package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	ID           uuid.UUID `json:"id"`
	CustomerName string    `json:"customer_name"`
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
