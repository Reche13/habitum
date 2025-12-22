package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerUserRoutes(users *echo.Group, h *handler.Handlers) {
		users.POST("", h.User.CreateUser)
		users.GET("", h.User.GetUsers)
		users.GET("/:id", h.User.GetUser)
		users.PUT("/:id", h.User.UpdateUser)
		users.DELETE("/:id", h.User.DeleteUser)
}