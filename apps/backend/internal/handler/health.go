package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/server"
)

type HealthHandler struct {
	Handler
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		Handler: NewHandler(s),
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	response := map[string]interface{}{
		"status":      "healthy",
		"timestamp":   time.Now().UTC(),
	}

	err := c.JSON(http.StatusOK, response)
	if err != nil {
		return fmt.Errorf("failed to write JSON response: %w", err)
	}

	return nil
}