package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerDashboardRoutes(dashboard *echo.Group, h *handler.Handlers) {
	dashboard.GET("/home", h.Dashboard.GetHome)
}


