package main

import (
	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/logger"
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

	srv := server.New(&cfg.Server, log)

	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}