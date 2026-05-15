package repository

import (
	"context"

	"order-service/features/order/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateTx(
	ctx context.Context,
	tx *gorm.DB,
	order *entity.Order,
) error {
	return tx.WithContext(ctx).
		Create(order).
		Error
}

func (r *orderRepository) CreateOutboxTx(
	ctx context.Context,
	tx *gorm.DB,
	outbox *entity.OutboxEvent,
) error {
	return tx.WithContext(ctx).
		Create(outbox).
		Error
}

func (r *orderRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Order, error) {
	var order entity.Order

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) FindAll(
	ctx context.Context,
) ([]entity.Order, error) {
	var orders []entity.Order

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&orders).
		Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
