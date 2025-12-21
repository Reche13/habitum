package handler

import "github.com/reche13/habitum/internal/server"

type Handlers struct {
	Health *HealthHandler
}

func NewHandlers(s *server.Server) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(s),
	}
}
