package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		HandleError: true,
		Skipper:     nil,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			requestID := GetRequestID(c)

			event := logger.Info().
				Str("request_id", requestID).
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("remote_ip", v.RemoteIP).
				Dur("latency", v.Latency).
				Str("latency_human", v.Latency.String())

			if v.Error != nil {
				event = event.Err(v.Error)
			}

			if v.Status >= 500 {
				event.Msg("HTTP request error")
			} else if v.Status >= 400 {
				event.Msg("HTTP request client error")
			} else {
				event.Msg("HTTP request")
			}

			return nil
		},
	})
}
