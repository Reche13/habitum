package handler

import (
	"github.com/reche13/habitum/internal/server"
	"github.com/reche13/habitum/internal/service"
)

type Handlers struct {
	Health *HealthHandler
	User *UserHandler
}

func NewHandlers(s *server.Server, service *service.UserService) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(s),
		User: NewUserHandler(s, service),
	}
}
