package dto

type CreateOrderRequest struct {
	CustomerName string  `json:"customer_name" validate:"required,min=3,max=100"`
	ProductName  string  `json:"product_name" validate:"required,min=2,max=100"`
	Quantity     int     `json:"quantity" validate:"required,gt=0"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
}
