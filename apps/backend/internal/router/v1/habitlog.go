package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerHabitLogRoutes(logs *echo.Group, h *handler.Handlers) {
	logs.POST("", h.HabitLog.Create)    
	logs.GET("", h.HabitLog.GetByDateRange)  
}