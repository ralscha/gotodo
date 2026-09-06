package request

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantError string
	}{
		{name: "valid", body: `{"name":"todo"}`},
		{name: "empty", body: "", wantError: "body must not be empty"},
		{name: "malformed", body: `{"name":`, wantError: "body contains badly-formed JSON"},
		{name: "unknown field", body: `{"unknown":true}`, wantError: `body contains unknown key "unknown"`},
		{name: "wrong type", body: `{"name":42}`, wantError: `incorrect JSON type for field "name"`},
		{name: "multiple values", body: `{"name":"one"} {"name":"two"}`, wantError: "single JSON value"},
		{
			name:      "too large",
			body:      `{"name":"` + strings.Repeat("a", maxBytes) + `"}`,
			wantError: "body must not be larger than 1048576 bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			var dst struct {
				Name string `json:"name"`
			}

			err := DecodeJSON(recorder, request, &dst)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("DecodeJSON() error = %v", err)
				}
				if dst.Name != "todo" {
					t.Fatalf("DecodeJSON() name = %q, want todo", dst.Name)
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("DecodeJSON() error = %v, want containing %q", err, tt.wantError)
			}
		})
	}
}

func TestDecodeJSONRejectsUnsupportedMediaType(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"todo"}`))
	request.Header.Set("Content-Type", "text/plain")

	err := DecodeJSON(recorder, request, &struct{}{})
	if !errors.Is(err, errUnsupportedMediaType) {
		t.Fatalf("DecodeJSON() error = %v, want unsupported media type", err)
	}
}
