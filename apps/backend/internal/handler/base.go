package handler

import (
	"github.com/reche13/habitum/internal/server"
)

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) Handler {
	return Handler{server: s}
}
