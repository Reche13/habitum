package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/rs/zerolog"
)

func ErrorHandler(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			httpErr, ok := err.(*errs.HTTPError)
			if ok {
				return c.JSON(httpErr.Status, httpErr)
			}

			if echoErr, ok := err.(*echo.HTTPError); ok {
				return c.JSON(echoErr.Code, map[string]interface{}{
					"code":    makeUpperCaseWithUnderscores(http.StatusText(echoErr.Code)),
					"message": echoErr.Message,
					"status":  echoErr.Code,
				})
			}

			requestID := c.Request().Header.Get("X-Request-ID")
			logger.Error().
				Err(err).
				Str("request_id", requestID).
				Str("path", c.Request().URL.Path).
				Str("method", c.Request().Method).
				Msg("unhandled error")

			return c.JSON(http.StatusInternalServerError, errs.NewInternalServerError("An unexpected error occurred"))
		}
	}
}

func makeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}

