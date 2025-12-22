package main

import (
	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/handler"
	"github.com/reche13/habitum/internal/logger"
	"github.com/reche13/habitum/internal/repository"
	"github.com/reche13/habitum/internal/router"
	"github.com/reche13/habitum/internal/server"
	"github.com/reche13/habitum/internal/service"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load config")
	}

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	repositories := repository.NewRepositories(srv.DB.Pool)
	services := service.NewServices(repositories)
	handlers := handler.NewHandlers(services)
	router := router.NewRouter(srv.Logger, handlers)

	srv.SetupHTTPServer(router)

	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
