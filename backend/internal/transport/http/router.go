package http

import (
	"github.com/labstack/echo/v4"
	"github.com/shaunpua/updoc/internal/doc"
)

func NewRouter(svc *doc.Service) *echo.Echo {
	e := echo.New()

	// Simple JSON health ping
	e.GET("/health", func(c echo.Context) error { return c.String(200, "ok") })

	// v1 group for future middleware (e.g., auth, logging)
	v1 := e.Group("/v1")

	h := &DocHandler{Svc: svc}
	v1.GET("/docs/:id", h.GetDoc)
	v1.POST("/docs/:id", h.UpdateDoc)

	return e
}
