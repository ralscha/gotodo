package main

import (
	"strings"
	"testing"

	"github.com/aarondl/null/v8"
)

func TestValidatedTodo(t *testing.T) {
	tests := []struct {
		name        string
		todo        ValidatedTodo
		errorField  string
		wantInvalid bool
	}{
		{name: "valid", todo: ValidatedTodo{Subject: "subject", Description: null.StringFrom("description")}},
		{name: "missing subject", todo: ValidatedTodo{}, errorField: "subject", wantInvalid: true},
		{name: "blank subject", todo: ValidatedTodo{Subject: "   "}, errorField: "subject", wantInvalid: true},
		{
			name:        "subject too long",
			todo:        ValidatedTodo{Subject: strings.Repeat("a", 256)},
			errorField:  "subject",
			wantInvalid: true,
		},
		{
			name:        "description too long",
			todo:        ValidatedTodo{Subject: "subject", Description: null.StringFrom(strings.Repeat("a", 256))},
			errorField:  "description",
			wantInvalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.todo.Validate()
			if got := errors.HasAny(); got != tt.wantInvalid {
				t.Fatalf("Validate().HasAny() = %t, want %t: %#v", got, tt.wantInvalid, errors.Errors)
			}
			if tt.errorField != "" && len(errors.Errors[tt.errorField]) == 0 {
				t.Fatalf("Validate() has no error for %q: %#v", tt.errorField, errors.Errors)
			}
		})
	}
}
