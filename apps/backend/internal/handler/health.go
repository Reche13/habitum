package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/server"
)

type HealthHandler struct {
	server *server.Server
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{server: s}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
	})
}
