package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
)

var validate = validator.New()

func Validate(payload interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := c.Bind(payload); err != nil {
				return errs.NewBadRequestError("Invalid request payload")
			}

			if err := validate.Struct(payload); err != nil {
				fieldErrors := make([]errs.FieldError, 0)
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, validationError := range validationErrors {
						fieldErrors = append(fieldErrors, errs.FieldError{
							Field: validationError.Field(),
							Error: getValidationErrorMessage(validationError),
						})
					}
				}
				return errs.NewValidationError(fieldErrors)
			}

			c.Set("validated_payload", payload)

			return next(c)
		}
	}
}

func ValidateStruct(s interface{}) []errs.FieldError {
	if err := validate.Struct(s); err != nil {
		fieldErrors := make([]errs.FieldError, 0)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationError := range validationErrors {
				fieldErrors = append(fieldErrors, errs.FieldError{
					Field: validationError.Field(),
					Error: getValidationErrorMessage(validationError),
				})
			}
		}
		return fieldErrors
	}
	return nil
}

func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	case "uuid":
		return err.Field() + " must be a valid UUID"
	default:
		return err.Field() + " is invalid"
	}
}

