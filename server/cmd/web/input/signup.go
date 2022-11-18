package input

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type SignUpInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *SignUpInput) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringLengthInRange{
			Name:    "password",
			Field:   s.Password,
			Message: "gte",
			Min:     8,
		},
		&validators.EmailIsPresent{
			Name:    "email",
			Field:   s.Email,
			Message: "email",
		},
	)
}
