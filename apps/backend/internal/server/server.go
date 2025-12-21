package server

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/config"
	"github.com/rs/zerolog"
)

type Server struct {
	log zerolog.Logger
	echo *echo.Echo
	port string
}

func New(cfg *config.ServerConfig, log zerolog.Logger) *Server {
	e := echo.New()

	return &Server{
		log: log,
		echo: e,
		port: cfg.Port,
	}
}

func (s *Server) Start() error {
	s.log.Info().
		Str("port", s.port).
		Msg("starting HTTP server")

	return s.echo.Start(":" + s.port)
}