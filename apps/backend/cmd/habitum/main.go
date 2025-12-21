package main

import (
	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load config")
	}

	log.Info().
		Str("port", cfg.Server.Port).
		Msg("starting habitum backend")

	log.Debug().
		Msg("debug logging is enabled")
}