package router

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func NewRouter(handler *handler.Handlers) *echo.Echo {
	router := echo.New()

	registerSystemRoutes(router, handler)

	router.POST("/users", handler.User.CreateUser)
	router.GET("/users", handler.User.GetUsers)
	
	return router
}