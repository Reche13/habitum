package handler

import (
	"strconv"

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
	
	breakdown, err := h.analyticsService.GetCategoryBreakdown(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, breakdown)
}

func (h *AnalyticsHandler) GetDayOfWeekAnalysis(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	
	// Get period query param (optional)
	period := c.QueryParam("period")
	
	analysis, err := h.analyticsService.GetDayOfWeekAnalysis(c.Request().Context(), userID, period)
	if err != nil {
		return err
	}

	return c.JSON(200, analysis)
}

func (h *AnalyticsHandler) GetMetrics(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	
	metrics, err := h.analyticsService.GetMetrics(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, metrics)
}

func (h *AnalyticsHandler) GetTopHabits(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	
	// Get query params
	limitStr := c.QueryParam("limit")
	limit := 10 // default
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	sortBy := c.QueryParam("sortBy")
	if sortBy == "" {
		sortBy = "completion" // default
	}

	topHabits, err := h.analyticsService.GetTopHabits(c.Request().Context(), userID, limit, sortBy)
	if err != nil {
		return err
	}

	return c.JSON(200, topHabits)
}

func (h *AnalyticsHandler) GetStreakLeaderboard(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	
	// Get limit query param
	limitStr := c.QueryParam("limit")
	limit := 10 // default
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	leaderboard, err := h.analyticsService.GetStreakLeaderboard(c.Request().Context(), userID, limit)
	if err != nil {
		return err
	}

	return c.JSON(200, leaderboard)
}

func (h *AnalyticsHandler) GetInsights(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")
	
	insights, err := h.analyticsService.GetInsights(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, insights)
}

