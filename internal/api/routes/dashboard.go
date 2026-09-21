package routes

import (
	"net/http"

	"github.com/dv1x3r/chinook-go/internal/api/routes/templates"
	"github.com/dv1x3r/chinook-go/internal/lib/echotempl"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetDashboard(c *echo.Context) error {
	return echotempl.Render(c, http.StatusOK, templates.Dashboard())
}
