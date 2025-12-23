package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerHabitRoutes(habits *echo.Group, h *handler.Handlers) {
		habits.POST("", h.Habit.CreateHabit)
		habits.GET("", h.Habit.GetHabits)
}