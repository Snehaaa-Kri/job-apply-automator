package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validate is your app's single validator instance.
// Node analogy: const schema = joi.object({ ... }) — reused across files.
var Validate = validator.New()

// FormatErrors takes raw validation errors and returns
// human-readable messages.
// Node analogy: converting Joi ValidationError into { field: message } map.
func FormatErrors(err error) map[string]string {
	errors := make(map[string]string)

	// Type-assert to ValidationErrors to iterate over each field error
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = formatTag(e)
	}

	return errors
}

// formatTag converts a validator tag into a readable message.
func formatTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
