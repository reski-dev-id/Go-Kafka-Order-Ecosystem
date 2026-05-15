package delivery

import (
	"net/http"

	"order-service/features/order/dto"
	"order-service/features/order/usecase"
	"order-service/internal/response"

	"github.com/google/uuid"
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

// CreateOrder godoc
// @Summary Create order
// @Description Create new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Create Order Request"
// @Success 201 {object} map[string]interface{} "Success create order"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/orders [post]
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

// GetOrders godoc
// @Summary Get orders
// @Description Get list orders with pagination
// @Tags Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit data"
// @Success 200 {object} map[string]interface{} "Success get orders"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetOrders(c echo.Context) error {

	page, limit := usecase.ParsePagination(
		c.QueryParam("page"),
		c.QueryParam("limit"),
	)

	result, total, err := h.usecase.GetOrders(
		c.Request().Context(),
		page,
		limit,
	)

	if err != nil {
		return err
	}

	return response.Success(
		c,
		http.StatusOK,
		"success get orders",
		map[string]interface{}{
			"items": dto.ToOrderResponses(result),
			"meta": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	)
}

// GetOrderByID godoc
// @Summary Get order by id
// @Description Get detail order by id
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Success get order"
// @Failure 400 {object} map[string]interface{} "Invalid order id"
// @Failure 404 {object} map[string]interface{} "Data not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c echo.Context) error {

	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid order id",
		)
	}

	result, err := h.usecase.GetOrderByID(
		c.Request().Context(),
		id,
	)

	if err != nil {
		return err
	}

	return response.Success(
		c,
		http.StatusOK,
		"success get order",
		dto.ToOrderResponse(result),
	)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update status order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Update Order Status Request"
// @Success 200 {object} map[string]interface{} "Success update order status"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Data not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {

	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid order id",
		)
	}

	var req dto.UpdateOrderStatusRequest

	err = c.Bind(&req)
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

	err = h.usecase.UpdateOrderStatus(
		c.Request().Context(),
		id,
		req,
	)

	if err != nil {
		return err
	}

	return response.Success(
		c,
		http.StatusOK,
		"success update order status",
		nil,
	)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check service health
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service healthy"
// @Router /health [get]
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
