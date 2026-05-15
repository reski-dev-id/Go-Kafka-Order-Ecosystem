package usecase

import (
	"context"
	"encoding/json"

	"order-service/features/order/dto"
	"order-service/features/order/entity"
	"order-service/features/order/repository"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderUsecase interface {
	CreateOrder(
		ctx context.Context,
		req dto.CreateOrderRequest,
	) (*entity.Order, error)

	GetOrders(
		ctx context.Context,
	) ([]entity.Order, error)
}

type orderUsecase struct {
	db        *gorm.DB
	orderRepo repository.OrderRepository
}

func NewOrderUsecase(
	db *gorm.DB,
	orderRepo repository.OrderRepository,
) OrderUsecase {
	return &orderUsecase{
		db:        db,
		orderRepo: orderRepo,
	}
}

func (u *orderUsecase) CreateOrder(
	ctx context.Context,
	req dto.CreateOrderRequest,
) (*entity.Order, error) {

	order := entity.Order{
		ID:           uuid.New(),
		CustomerName: req.CustomerName,
		ProductName:  req.ProductName,
		Quantity:     req.Quantity,
		Amount:       req.Amount,
		Status:       "PENDING",
	}

	payload, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	outbox := entity.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "order",
		AggregateID:   order.ID,
		EventType:     "order.created",
		Payload:       datatypes.JSON(payload),
		Status:        "PENDING",
	}

	tx := u.db.Begin()

	err = u.orderRepo.CreateTx(ctx, tx, &order)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = u.orderRepo.CreateOutboxTx(ctx, tx, &outbox)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (u *orderUsecase) GetOrders(
	ctx context.Context,
) ([]entity.Order, error) {
	return u.orderRepo.FindAll(ctx)
}
