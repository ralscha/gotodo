package input

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type EmailChangeInput struct {
	Password string `json:"password"`
	NewEmail string `json:"newEmail"`
}

func (e *EmailChangeInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringLengthInRange{
			Name:    "password",
			Field:   e.Password,
			Message: "gte",
			Min:     8,
		},
		maximumLength("password", e.Password, MaxPasswordLength),
		&validators.EmailIsPresent{
			Name:    "newEmail",
			Field:   e.NewEmail,
			Message: "email",
		},
		maximumLength("newEmail", e.NewEmail, MaxEmailLength),
	)
}

type PasswordInput struct {
	Password string `json:"password"`
}

func (p *PasswordInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "password",
			Field:   p.Password,
			Message: "required",
		},
		maximumLength("password", p.Password, MaxPasswordLength),
	)
}

type TokenInput struct {
	Token string `json:"token"`
}

func (t *TokenInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "token",
			Field:   t.Token,
			Message: "required",
		},
		maximumLength("token", t.Token, MaxTokenLength),
	)
}

type PasswordChangeInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (p *PasswordChangeInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "oldPassword",
			Field:   p.OldPassword,
			Message: "required",
		},
		maximumLength("oldPassword", p.OldPassword, MaxPasswordLength),
		&validators.StringLengthInRange{
			Name:    "newPassword",
			Field:   p.NewPassword,
			Message: "gte",
			Min:     8,
		},
		maximumLength("newPassword", p.NewPassword, MaxPasswordLength),
	)
}
