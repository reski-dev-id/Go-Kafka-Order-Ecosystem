package dto

type OrderCreatedEvent struct {
	EventID      string  `json:"EventID"`
	ID           string  `json:"ID"`
	CustomerName string  `json:"CustomerName"`
	ProductName  string  `json:"ProductName"`
	Quantity     int     `json:"Quantity"`
	Amount       float64 `json:"Amount"`
	Status       string  `json:"Status"`
	CreatedAt    string  `json:"CreatedAt"`
	UpdatedAt    string  `json:"UpdatedAt"`
}
