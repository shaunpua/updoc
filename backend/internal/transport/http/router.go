package http

import (
	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	// Simple JSON health ping
	e.GET("/health", func(c echo.Context) error { return c.String(200, "ok") })

	return e
}
