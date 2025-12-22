package handler

import "github.com/reche13/habitum/internal/service"

type Handlers struct {
	Health *HealthHandler
	User   *UserHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(),
		User:   NewUserHandler(services.User),
	}
}
