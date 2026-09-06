package response

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gobuffalo/validate"
)

func errorMessage(w http.ResponseWriter, status int, message string, headers http.Header) {
	JSONWithHeaders(w, status, map[string]string{"error": message}, headers)
}

func InternalServerError(w http.ResponseWriter, err error) {
	slog.Error(err.Error(), "error", err)

	message := "The server encountered a problem and could not process your request"
	errorMessage(w, http.StatusInternalServerError, message, nil)
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	message := "The requested resource could not be found"
	errorMessage(w, http.StatusNotFound, message, nil)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	errorMessage(w, http.StatusMethodNotAllowed, message, nil)
}

func BadRequest(w http.ResponseWriter, err error) {
	errorMessage(w, http.StatusBadRequest, err.Error(), nil)
}

func UnsupportedMediaType(w http.ResponseWriter) {
	errorMessage(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json", nil)
}

func FailedValidation(w http.ResponseWriter, v *validate.Errors) {
	errors := make(map[string][]string, len(v.Errors))
	for field, messages := range v.Errors {
		errors[lowerCamelCase(field)] = messages
	}
	JSON(w, http.StatusUnprocessableEntity, &validate.Errors{Errors: errors})
}

func lowerCamelCase(value string) string {
	parts := strings.Split(value, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func Unauthorized(w http.ResponseWriter) {
	message := "You must be authenticated to access this resource"
	errorMessage(w, http.StatusUnauthorized, message, nil)
}
