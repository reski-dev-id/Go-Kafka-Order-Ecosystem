package usecase

import (
	"context"
	"encoding/json"
	"strconv"

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
		page int,
		limit int,
	) ([]entity.Order, int64, error)

	GetOrderByID(
		ctx context.Context,
		id string,
	) (*entity.Order, error)

	UpdateOrderStatus(
		ctx context.Context,
		id string,
		req dto.UpdateOrderStatusRequest,
	) error
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
	page int,
	limit int,
) ([]entity.Order, int64, error) {

	offset := (page - 1) * limit

	orders, err := u.orderRepo.FindAll(
		ctx,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, err
	}

	total, err := u.orderRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (u *orderUsecase) GetOrderByID(
	ctx context.Context,
	id string,
) (*entity.Order, error) {

	orderID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.FindByID(
		ctx,
		orderID,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *orderUsecase) UpdateOrderStatus(
	ctx context.Context,
	id string,
	req dto.UpdateOrderStatusRequest,
) error {

	orderID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = u.orderRepo.FindByID(
		ctx,
		orderID,
	)

	if err != nil {
		return err
	}

	return u.orderRepo.UpdateStatus(
		ctx,
		orderID,
		req.Status,
	)
}

func ParsePagination(
	pageStr string,
	limitStr string,
) (int, int) {

	page := 1
	limit := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	return page, limit
}
