package delivery

import (
	"net/http"

	"order-service/features/order/dto"
	"order-service/features/order/usecase"
	"order-service/internal/response"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(
	usecase usecase.OrderUsecase,
) *OrderHandler {
	return &OrderHandler{
		usecase: usecase,
	}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req dto.CreateOrderRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	err = c.Validate(&req)
	if err != nil {
		return err
	}

	result, err := h.usecase.CreateOrder(
		c.Request().Context(),
		req,
	)

	if err != nil {
		return err
	}

	return response.Success(
		c,
		http.StatusCreated,
		"success create order",
		dto.ToOrderResponse(result),
	)
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	result, err := h.usecase.GetOrders(
		c.Request().Context(),
	)

	if err != nil {
		return err
	}

	return response.Success(
		c,
		http.StatusOK,
		"success get orders",
		dto.ToOrderResponses(result),
	)
}

func (h *OrderHandler) HealthCheck(c echo.Context) error {
	return response.Success(
		c,
		http.StatusOK,
		"service is healthy",
		map[string]interface{}{
			"service": "order-service",
			"status":  "UP",
		},
	)
}
