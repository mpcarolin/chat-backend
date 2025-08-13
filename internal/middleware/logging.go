package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			info := fmt.Sprintf("[%s] %s %s", time.Now().UTC().String(), req.Method, req.URL.Path)
			slog.Info(info)

			return next(c)
		}
	}
}
