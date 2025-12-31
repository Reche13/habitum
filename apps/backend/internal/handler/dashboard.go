package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/service"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetHome(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	dashboard, err := h.dashboardService.GetHome(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, dashboard)
}


