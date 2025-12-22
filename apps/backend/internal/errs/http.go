package errs

import (
	"net/http"
)

func NewUnauthorizedError(message string) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusForbidden)),
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NewBadRequestError(message string) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest)),
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewBadRequestErrorWithFields(message string, fields []FieldError) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest)),
		Message: message,
		Status:  http.StatusBadRequest,
		Fields:  fields,
	}
}

func NewNotFoundError(message string) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusNotFound)),
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func NewInternalServerError(message string) *HTTPError {
	if message == "" {
		message = http.StatusText(http.StatusInternalServerError)
	}
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError)),
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewValidationError(fields []FieldError) *HTTPError {
	return &HTTPError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Status:  http.StatusBadRequest,
		Fields:  fields,
	}
}
