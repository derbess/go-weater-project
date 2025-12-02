package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: getErrorMessage(e),
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
			})
		}
	}

	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed validation on tag '%s'", e.Field(), e.Tag())
	}
}

func RespondWithValidationError(w http.ResponseWriter, err error) {
	errors := FormatValidationErrors(err)

	response := ValidationErrorResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}




