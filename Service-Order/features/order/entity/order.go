package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	CustomerName string
	ProductName  string
	Quantity     int
	Amount       float64
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Order) TableName() string {
	return "orders"
}
