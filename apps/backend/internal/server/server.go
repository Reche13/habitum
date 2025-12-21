package server

import (
	"errors"
	"net/http"

	"github.com/reche13/habitum/internal/config"
	"github.com/rs/zerolog"
)

type Server struct {
	Config *config.Config
	Logger zerolog.Logger
	httpServer *http.Server
}

func New(cfg *config.Config, logger zerolog.Logger) *Server {
	return &Server{
		Logger: logger,
		Config: cfg,
	}
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