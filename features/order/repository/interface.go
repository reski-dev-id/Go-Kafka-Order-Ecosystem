package repository

import (
	"context"

	"order-service/features/order/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateTx(
		ctx context.Context,
		tx *gorm.DB,
		order *entity.Order,
	) error

	CreateOutboxTx(
		ctx context.Context,
		tx *gorm.DB,
		outbox *entity.OutboxEvent,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.Order, error)

	FindAll(
		ctx context.Context,
	) ([]entity.Order, error)
}
