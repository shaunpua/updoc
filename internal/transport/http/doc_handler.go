package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shaunpua/updoc/internal/doc"
	"github.com/shaunpua/updoc/internal/providers/confluence"
)

type DocHandler struct {
	Svc interface {
		GetPage(string) (*confluence.PageResp, error)
		UpdatePage(pageID, html, status string) (doc.DocFlag, error)
		ListFlags(pageID string) []doc.DocFlag
	}
}

func (h *DocHandler) GetDoc(c echo.Context) error {
	id := c.Param("id")
	page, err := h.Svc.GetPage(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	// inline flags for convenience
	flags := h.Svc.ListFlags(id)
	return c.JSON(http.StatusOK, echo.Map{
		"page":  page,
		"flags": flags,
	})
}

type updateReq struct {
	Body   string `json:"body,omitempty"`
	Status string `json:"status" validate:"required"`
}

func (h *DocHandler) UpdateDoc(c echo.Context) error {
	id := c.Param("id")
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	flag, err := h.Svc.UpdatePage(id, req.Body, req.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	return c.JSON(http.StatusOK, flag)
}
