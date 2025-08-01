package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

var validate = validator.New()

// readJSON decodes JSON into the given destination and validates it using go-playground/validator.
func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}, customMessages map[string]string) error {
	maxBytes := 1_048_576 // Set a limit for the request body size.
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Decode the JSON request body into the destination struct.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return handleJSONDecodeError(err)
	}

	// Validate the decoded struct using go-playground/validator.
	err = validate.Struct(dst)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// If no custom messages provided, use empty map
			if customMessages == nil {
				customMessages = map[string]string{}
			}
			return handleValidationErrors(validationErrors, customMessages)
		}
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// handleJSONDecodeError provides custom error messages for JSON decoding issues.
func handleJSONDecodeError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("body contains unknown field %s", fieldName)
	case err.Error() == "http: request body too large":
		return fmt.Errorf("body must not be larger than %d bytes", 1_048_576)
	default:
		return err
	}
}

// handleValidationErrors converts validator.ValidationErrors into a user-friendly error message.
func handleValidationErrors(validationErrors validator.ValidationErrors, customMessages map[string]string) error {
	var errorMessages []string
	for _, fieldError := range validationErrors {
		field := fieldError.Field()
		tag := fieldError.Tag()
		key := fmt.Sprintf("%s.%s", field, tag)

		// Check if a custom message exists for the field and tag.
		if msg, exists := customMessages[key]; exists {
			errorMessages = append(errorMessages, msg)
		} else {
			// Fallback to default error message if no custom message is defined.
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation on tag '%s'", field, tag))
		}
	}
	return errors.New(strings.Join(errorMessages, ", "))
}
