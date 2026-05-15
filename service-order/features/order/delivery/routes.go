package delivery

import (
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	handler *OrderHandler,
) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", handler.HealthCheck)

	api := e.Group("/api/v1")

	api.POST("/orders", handler.CreateOrder)
	api.GET("/orders", handler.GetOrders)
	api.GET("/orders/:id", handler.GetOrderByID)
	api.PATCH("/orders/:id/status", handler.UpdateOrderStatus)
}
