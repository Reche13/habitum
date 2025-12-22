package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/database"
	"github.com/rs/zerolog"
)

type Server struct {
	Config     *config.Config
	Logger     zerolog.Logger
	DB         *database.Database
	httpServer *http.Server
}

func New(cfg *config.Config, logger zerolog.Logger) (*Server, error) {
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr: ":" + s.Config.Server.Port,
		Handler: handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().Str("port", s.Config.Server.Port).Msg("starting server")


	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		s.Logger.Info().Str("signal", sig.String()).Msg("shutting down server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.Logger.Info().Msg("server exited gracefully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.DB != nil {
		s.DB.Close()
	}
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}