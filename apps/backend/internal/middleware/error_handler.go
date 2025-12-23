package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/sqlerr"
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
				return c.JSON(httpErr.Status, addRequestIDToError(httpErr, c))
			}

			if echoErr, ok := err.(*echo.HTTPError); ok {
				return c.JSON(echoErr.Code, map[string]interface{}{
					"code":      errs.MakeUpperCaseWithUnderscores(http.StatusText(echoErr.Code)),
					"message":   echoErr.Message,
					"status":    echoErr.Code,
					"request_id": GetRequestID(c),
				})
			}

			if sqlerr.IsDatabaseError(err) {
				var resourceName string

				var dbErr *sqlerr.DatabaseError
				if errors.As(err, &dbErr) {
					resourceName = dbErr.ResourceName
					err = dbErr.Unwrap()
				} else {
					resourceName = "resource"
				}
				
				httpErr := sqlerr.HandleError(err, resourceName)
				if httpErr != nil {
					if convertedErr, ok := httpErr.(*errs.HTTPError); ok {
						return c.JSON(convertedErr.Status, addRequestIDToError(convertedErr, c))
					}
				}
			}

			requestID := GetRequestID(c)
			logger.Error().
				Err(err).
				Str("request_id", requestID).
				Str("path", c.Request().URL.Path).
				Str("method", c.Request().Method).
				Msg("unhandled error")

			return c.JSON(http.StatusInternalServerError, addRequestIDToError(
				errs.NewInternalServerError("An unexpected error occurred"),
				c,
			))
		}
	}
}


func getHTTPErrorStatus(err error) int {
	if httpErr, ok := err.(*errs.HTTPError); ok {
		return httpErr.Status
	}
	return http.StatusInternalServerError
}

func addRequestIDToError(httpErr *errs.HTTPError, c echo.Context) map[string]interface{} {
	response := map[string]interface{}{
		"code":      httpErr.Code,
		"message":   httpErr.Message,
		"status":    httpErr.Status,
		"request_id": GetRequestID(c),
	}
	
	if len(httpErr.Fields) > 0 {
		response["fields"] = httpErr.Fields
	}
	
	return response
}

