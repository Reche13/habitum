package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerCalendarRoutes(calendar *echo.Group, h *handler.Handlers) {
	calendar.GET("/completions", h.Calendar.GetCompletions)
	calendar.GET("/month", h.Calendar.GetMonth)
	calendar.GET("/week", h.Calendar.GetWeek)
	calendar.GET("/year", h.Calendar.GetYear)
}

