package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shaunpua/updoc/internal/services"
)

type OrganizationHandler struct {
	orgService        *services.OrganizationService
	confluenceService *services.ConfluenceService
}

func NewOrganizationHandler(orgService *services.OrganizationService, confluenceService *services.ConfluenceService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService:        orgService,
		confluenceService: confluenceService,
	}
}

// CreateOrganization handles POST /api/v1/orgs
func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	var req services.CreateOrgRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Basic validation
	if req.Name == "" || req.UserEmail == "" || req.UserName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name, user_name, and user_email are required")
	}

	resp, err := h.orgService.CreateWithUser(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

// GetOrganization handles GET /api/v1/orgs/:slug
func (h *OrganizationHandler) GetOrganization(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "slug is required")
	}

	org, err := h.orgService.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, org)
}

// TestConfluence handles POST /api/v1/orgs/:id/test-confluence
func (h *OrganizationHandler) TestConfluence(c echo.Context) error {
	orgID := c.Param("id")
	if orgID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "organization ID is required")
	}

	result, err := h.confluenceService.TestConnection(c.Request().Context(), orgID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// ListConfluencePages handles GET /api/v1/orgs/:id/confluence/pages
func (h *OrganizationHandler) ListConfluencePages(c echo.Context) error {
	orgID := c.Param("id")
	if orgID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "organization ID is required")
	}

	limit := 10 // default
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	pages, err := h.confluenceService.ListPages(c.Request().Context(), orgID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"pages": pages,
		"count": len(pages),
	})
}
