package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/database"
	"github.com/rs/zerolog"
)

type Server struct {
	Config *config.Config
	Logger zerolog.Logger
	DB *database.Database
	httpServer *http.Server
}

func New(cfg *config.Config, logger zerolog.Logger)(*Server, error) {
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		Config: cfg,
		Logger: logger,
		DB: db,
	}, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr: ":" + s.Config.Server.Port,
		Handler: handler,
	}
}


func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().Str("post", s.Config.Server.Port).Msg("starting server")

	return s.httpServer.ListenAndServe()
}