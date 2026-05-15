package middleware

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func CustomErrorHandler(err error, c echo.Context) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		var validationErrors []map[string]string

		for _, fe := range ve {
			validationErrors = append(
				validationErrors,
				map[string]string{
					"field": fe.Field(),
					"tag":   fe.Tag(),
				},
			)
		}

		_ = c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Success: false,
				Message: "validation error",
				Errors:  validationErrors,
			},
		)

		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = c.JSON(
			http.StatusNotFound,
			ErrorResponse{
				Success: false,
				Message: "data not found",
			},
		)

		return
	}

	var he *echo.HTTPError

	if errors.As(err, &he) {
		_ = c.JSON(
			he.Code,
			ErrorResponse{
				Success: false,
				Message: he.Message.(string),
			},
		)

		return
	}

	_ = c.JSON(
		http.StatusInternalServerError,
		ErrorResponse{
			Success: false,
			Message: "internal server error",
		},
	)
}
