package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/service"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}


func (h *AnalyticsHandler) GetCompletionTrend(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	
	// Get period query param (default to "30d")
	period := c.QueryParam("period")
	if period == "" {
		period = "30d"
	}

	// Validate period
	validPeriods := map[string]bool{"7d": true, "30d": true, "90d": true, "all": true}
	if !validPeriods[period] {
		return errs.NewBadRequestError("Invalid period. Must be one of: 7d, 30d, 90d, all")
	}

	trend, err := h.analyticsService.GetCompletionTrend(c.Request().Context(), userID, period)
	if err != nil {
		return err
	}

	return c.JSON(200, trend)
}

func (h *AnalyticsHandler) GetCategoryBreakdown(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	_ = userID // Placeholder
	return c.JSON(200, map[string]interface{}{"message": "Not implemented yet"})
}

func (h *AnalyticsHandler) GetDayOfWeekAnalysis(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	_ = userID // Placeholder
	return c.JSON(200, map[string]interface{}{"message": "Not implemented yet"})
}

func (h *AnalyticsHandler) GetMetrics(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	_ = userID // Placeholder
	return c.JSON(200, map[string]interface{}{"message": "Not implemented yet"})
}

func (h *AnalyticsHandler) GetTopHabits(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	_ = userID // Placeholder
	return c.JSON(200, map[string]interface{}{"message": "Not implemented yet"})
}

func (h *AnalyticsHandler) GetStreakLeaderboard(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	_ = userID // Placeholder
	return c.JSON(200, map[string]interface{}{"message": "Not implemented yet"})
}

func (h *AnalyticsHandler) GetInsights(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	_ = userID // Placeholder
	return c.JSON(200, map[string]interface{}{"message": "Not implemented yet"})
}

