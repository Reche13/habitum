package router

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerSystemRoutes(e *echo.Echo, h *handler.Handlers) {
	e.GET("/health", h.Health.CheckHealth)
	e.GET("/status", h.Health.CheckHealth)
}