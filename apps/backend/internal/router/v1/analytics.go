package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerAnalyticsRoutes(analytics *echo.Group, h *handler.Handlers) {
	analytics.GET("/completion-trend", h.Analytics.GetCompletionTrend)
	analytics.GET("/category-breakdown", h.Analytics.GetCategoryBreakdown)
	analytics.GET("/day-of-week", h.Analytics.GetDayOfWeekAnalysis)
	analytics.GET("/metrics", h.Analytics.GetMetrics)
	analytics.GET("/top-habits", h.Analytics.GetTopHabits)
	analytics.GET("/streak-leaderboard", h.Analytics.GetStreakLeaderboard)
	analytics.GET("/insights", h.Analytics.GetInsights)
}