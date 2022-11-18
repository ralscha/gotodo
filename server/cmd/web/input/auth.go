package input

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Validatable interface {
	Validate() *validate.Errors
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
		&validators.EmailIsPresent{
			Name:    "email",
			Field:   l.Email,
			Message: "email",
		},
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
		&validators.StringIsPresent{
			Name:    "resetToken",
			Field:   p.ResetToken,
			Message: "required",
		},
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
	)
}
