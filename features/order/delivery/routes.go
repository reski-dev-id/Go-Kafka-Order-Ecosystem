package delivery

import "github.com/labstack/echo/v4"

func RegisterRoutes(
	e *echo.Echo,
	handler *OrderHandler,
) {
	e.GET("/health", handler.HealthCheck)

	api := e.Group("/api/v1")

	api.POST("/orders", handler.CreateOrder)
	api.GET("/orders", handler.GetOrders)
}
