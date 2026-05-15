package response

import "github.com/labstack/echo/v4"

func Success(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Error(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]interface{}{
		"success": false,
		"message": message,
	})
}
