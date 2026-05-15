package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		stop := time.Now()

		latency := stop.Sub(start)

		req := c.Request()
		res := c.Response()

		log.Printf(
			"%s %s %d %s %s",
			req.Method,
			req.RequestURI,
			res.Status,
			latency,
			c.RealIP(),
		)

		return err
	}
}
