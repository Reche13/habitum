package main

import (
	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/handler"
	"github.com/reche13/habitum/internal/logger"
	"github.com/reche13/habitum/internal/router"
	"github.com/reche13/habitum/internal/server"
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

	h := handler.NewHandlers(srv)
	r := router.NewRouter(h)

	srv.SetupHTTPServer(r)

	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}