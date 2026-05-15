package bootstrap

import (
	orderDelivery "order-service/features/order/delivery"
	orderRepository "order-service/features/order/repository"
	orderUsecase "order-service/features/order/usecase"

	"order-service/internal/config"
	"order-service/internal/database"
	customMiddleware "order-service/internal/middleware"
	customValidator "order-service/internal/pkg/validator"

	"github.com/labstack/echo/v4"
)

func NewApp() (*echo.Echo, *config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		return nil, nil, err
	}

	orderRepo := orderRepository.NewOrderRepository(db)

	orderUC := orderUsecase.NewOrderUsecase(
		db,
		orderRepo,
	)

	orderHandler := orderDelivery.NewOrderHandler(
		orderUC,
	)

	e := echo.New()

	e.Validator = customValidator.NewValidator()

	e.HTTPErrorHandler = customMiddleware.CustomErrorHandler

	e.Use(customMiddleware.RequestLogger)

	orderDelivery.RegisterRoutes(
		e,
		orderHandler,
	)

	return e, cfg, nil
}
