package input

import (
	"strings"
	"testing"
)

func TestInputLengthLimits(t *testing.T) {
	tests := []struct {
		name  string
		input Validatable
		field string
	}{
		{
			name:  "signup email",
			input: &SignUpInput{Email: strings.Repeat("a", MaxEmailLength+1), Password: "valid-password"},
			field: "email",
		},
		{
			name:  "signup password",
			input: &SignUpInput{Email: "person@example.com", Password: strings.Repeat("a", MaxPasswordLength+1)},
			field: "password",
		},
		{
			name:  "reset token",
			input: &PasswordResetInput{Password: "valid-password", ResetToken: strings.Repeat("a", MaxTokenLength+1)},
			field: "reset_token",
		},
		{
			name:  "new email",
			input: &EmailChangeInput{Password: "valid-password", NewEmail: strings.Repeat("a", MaxEmailLength+1)},
			field: "new_email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.Validate()
			if len(errors.Errors[tt.field]) == 0 {
				t.Fatalf("Validate() has no error for %q: %#v", tt.field, errors.Errors)
			}
		})
	}
}
