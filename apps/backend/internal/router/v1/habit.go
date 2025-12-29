package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerHabitRoutes(habits *echo.Group, h *handler.Handlers) {
	habits.POST("", h.Habit.CreateHabit)
	habits.GET("", h.Habit.GetHabits)
	habits.GET("/:id", h.Habit.GetHabit)
	habits.PATCH("/:id", h.Habit.UpdateHabit)
	habits.DELETE("/:id", h.Habit.DeleteHabit)

	// Completion endpoints
	habits.POST("/:id/complete", h.Habit.MarkComplete)
	habits.DELETE("/:id/complete", h.Habit.UnmarkComplete)
	habits.GET("/:id/completions", h.Habit.GetCompletions)
	habits.GET("/:id/completion-history", h.Habit.GetCompletionHistory)

	habitLogs := habits.Group("/:habit_id/logs")
	registerHabitLogRoutes(habitLogs, h)
}