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
		&validators.EmailIsPresent{
			Name:    "newEmail",
			Field:   e.NewEmail,
			Message: "email",
		},
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
		&validators.StringLengthInRange{
			Name:    "newPassword",
			Field:   p.NewPassword,
			Message: "gte",
			Min:     8,
		},
	)
}
