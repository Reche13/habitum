package handler

import "github.com/reche13/habitum/internal/service"

type Handlers struct {
	Health *HealthHandler
	User   *UserHandler
	Habit *HabitHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(),
		User:   NewUserHandler(services.User),
		Habit: NewHabitHandler(services.Habit),
	}
}
