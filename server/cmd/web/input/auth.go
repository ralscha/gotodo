package input

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Validatable interface {
	Validate() *validate.Errors
}

const (
	MaxEmailLength    = 254
	MaxPasswordLength = 128
	MaxTokenLength    = 64
)

func maximumLength(name, field string, max int) validate.Validator {
	return &validators.StringLengthInRange{
		Name:    name,
		Field:   field,
		Message: "lte",
		Min:     0,
		Max:     max,
	}
}

type LoginInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (l *LoginInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringLengthInRange{
			Name:    "password",
			Field:   l.Password,
			Message: "gte",
			Min:     8,
		},
		maximumLength("password", l.Password, MaxPasswordLength),
		&validators.EmailIsPresent{
			Name:    "email",
			Field:   l.Email,
			Message: "email",
		},
		maximumLength("email", l.Email, MaxEmailLength),
	)
}

type PasswordResetInput struct {
	Password   string `json:"password"`
	ResetToken string `json:"resetToken"`
}

func (p *PasswordResetInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringLengthInRange{
			Name:    "password",
			Field:   p.Password,
			Message: "gte",
			Min:     8,
		},
		maximumLength("password", p.Password, MaxPasswordLength),
		&validators.StringIsPresent{
			Name:    "resetToken",
			Field:   p.ResetToken,
			Message: "required",
		},
		maximumLength("resetToken", p.ResetToken, MaxTokenLength),
	)
}

type PasswordResetRequestInput struct {
	Email string `json:"email"`
}

func (p *PasswordResetRequestInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.EmailIsPresent{
			Name:    "email",
			Field:   p.Email,
			Message: "email",
		},
		maximumLength("email", p.Email, MaxEmailLength),
	)
}
