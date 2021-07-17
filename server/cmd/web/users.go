package main

import (
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotodo.rasc.ch/internal/models"
	"net/http"
	"time"
)

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,gte=8"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	valid, fieldErrors := app.validate(input)
	if !valid {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, SaveResponse{
			Success:     false,
			FieldErrors: fieldErrors,
		})
		return
	}

	hash, err := argon2id.CreateHash(input.Password, &argon2id.Params{
		Memory:      app.config.Argon2.Memory,
		Iterations:  app.config.Argon2.Iterations,
		Parallelism: app.config.Argon2.Parallelism,
		SaltLength:  app.config.Argon2.SaltLength,
		KeyLength:   app.config.Argon2.KeyLength,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newUser := models.AppUser{
		Email:        input.Email,
		PasswordHash: hash,
		Authority:    "user",
		Activated:    false,
	}

	ctx, cancel := app.createDbContext()
	err = newUser.Insert(ctx, app.db, boil.Infer())
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.insertToken(newUser.ID, 24*time.Hour, scopeSignup)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"confirmationLink": token.plain,
		}

		err = app.mailer.Send(newUser.Email, "signup-confirm.tmpl", data)
		if err != nil {
			app.logger.Error(err)
		}
	})

	w.WriteHeader(http.StatusNoContent)
}

/*
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
*/
