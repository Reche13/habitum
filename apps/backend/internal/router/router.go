package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reche13/habitum/internal/handler"
	mw "github.com/reche13/habitum/internal/middleware"
	v1 "github.com/reche13/habitum/internal/router/v1"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger ,handlers *handler.Handlers) *echo.Echo {
	router := echo.New()
	
	router.Use(mw.Recover())
	router.Use(mw.RequestID())
	router.Use(mw.Logger(logger))
	router.Use(mw.CORS())
	router.Use(middleware.BodyLimit("2M"))
	router.Use(mw.ErrorHandler(logger))

	registerSystemRoutes(router, handlers)
	apiV1 := router.Group("/api/v1")
	v1.RegisterAPIV1Routes(apiV1, handlers)

	return router
}



