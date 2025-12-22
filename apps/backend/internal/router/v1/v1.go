package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func RegisterAPIV1Routes(api *echo.Group, h *handler.Handlers) {
	users := api.Group("/users")
	registerUserRoutes(users, h)
}
