package errs

import (
	"strings"
)

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type HTTPError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Status  int          `json:"status"`
	Fields  []FieldError `json:"fields,omitempty"` // Only for validation errors
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Is(target error) bool {
	_, ok := target.(*HTTPError)
	return ok
}

func (e *HTTPError) WithMessage(message string) *HTTPError {
	return &HTTPError{
		Code:    e.Code,
		Message: message,
		Status:  e.Status,
		Fields:  e.Fields,
	}
}

func MakeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}
