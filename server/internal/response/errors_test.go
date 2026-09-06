package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/validate"
)

func TestFailedValidationUsesJSONFieldNames(t *testing.T) {
	recorder := httptest.NewRecorder()
	errors := &validate.Errors{Errors: map[string][]string{
		"new_email": {"email"},
	}}

	FailedValidation(recorder, errors)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if got := recorder.Body.String(); got != `{"errors":{"newEmail":["email"]}}` {
		t.Fatalf("body = %q", got)
	}
}
