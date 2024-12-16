package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Custom validation struct
type CustomValidator struct {
	validator *validator.Validate
}

// Custom validate function
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)

		// Map custom messages for each field
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			switch fieldErr.Tag() {
			case "required":
				errorMessages[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				errorMessages[fieldName] = fmt.Sprintf("%s must be a valid email", fieldName)
			case "min":
				errorMessages[fieldName] = fmt.Sprintf("%s must be at least %s characters", fieldName, fieldErr.Param())
			default:
				errorMessages[fieldName] = fmt.Sprintf("%s is invalid", fieldName)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, errorMessages)
	}
	return nil
}
