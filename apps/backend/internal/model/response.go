package model

// APIResponse is a standard API response wrapper
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Error represents an API error
type Error struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

// FieldError represents a field-level validation error
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Meta contains pagination and other metadata
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Total     int    `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// SuccessResponse creates a successful API response
func SuccessResponse[T any](data T) *APIResponse[T] {
	return &APIResponse[T]{
		Success: true,
		Data:    data,
	}
}

// SuccessResponseWithMeta creates a successful API response with metadata
func SuccessResponseWithMeta[T any](data T, meta *Meta) *APIResponse[T] {
	return &APIResponse[T]{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(code, message string, fields []FieldError) *APIResponse[interface{}] {
	return &APIResponse[interface{}]{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	}
}


