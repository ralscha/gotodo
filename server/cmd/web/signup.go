package main

import (
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"net/http"
	"time"
)

func (app *application) passwordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	email, err := app.readString(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if email == "" {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{"email": "required"},
		})
		return
	}

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID),
		models.AppUserWhere.Email.EQ(email)).One(r.Context(), app.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.serverErrorResponse(w, r, err)
		return
	}

	if user != nil {
		token, err := app.insertToken(r.Context(), user.ID, app.config.Cleanup.PasswordResetTokenMaxAge, scopePasswordReset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		app.background(func() {
			data := map[string]interface{}{
				"resetLink": app.config.BaseUrl + "#/password-reset/" + token.plain,
			}

			err = app.mailer.Send(email, "password-reset.tmpl", data)
			if err != nil {
				app.logger.Error(err)
			}
		})
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) passwordResetHandler(w http.ResponseWriter, r *http.Request) {
	var resetInput struct {
		Password   string `name:"password" validate:"required,gte=8"`
		ResetToken string `name:"resetToken" validate:"required"`
	}
	if ok := app.parseFromForm(w, r, &resetInput); !ok {
		return
	}

	userId, err := app.getAppUserIdFromToken(r.Context(), scopePasswordReset, resetInput.ResetToken)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if userId == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	compromised, err := app.isPasswordCompromised(resetInput.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if compromised {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{
				"password": "weak",
			},
		})
		return
	}

	newPasswordHash, err := argon2id.CreateHash(resetInput.Password, &argon2id.Params{
		Memory:      app.config.Argon2.Memory,
		Iterations:  app.config.Argon2.Iterations,
		Parallelism: app.config.Argon2.Parallelism,
		SaltLength:  app.config.Argon2.SaltLength,
		KeyLength:   app.config.Argon2.KeyLength,
	})

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userId)).UpdateAll(r.Context(), app.db,
		models.M{models.AppUserColumns.PasswordHash: newPasswordHash})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.deleteAllTokensForUser(r.Context(), userId, scopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `name:"email" validate:"required,email"`
		Password string `name:"password" validate:"required,gte=8"`
	}
	if ok := app.parseFromForm(w, r, &input); !ok {
		return
	}

	exists, err := models.AppUsers(
		models.AppUserWhere.Email.EQ(input.Email),
		qm.Or2(models.AppUserWhere.EmailNew.EQ(null.NewString(input.Email, true))),
	).Exists(r.Context(), app.db)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if exists {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{
				"email": "exists",
			},
		})
		return
	}

	compromised, err := app.isPasswordCompromised(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if compromised {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{
				"password": "weak",
			},
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
		Authority:    "USER",
		Activated:    false,
	}

	err = newUser.Insert(r.Context(), app.db, boil.Infer())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.insertToken(r.Context(), newUser.ID, app.config.Cleanup.SignupTokenMaxAge, scopeSignup)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"confirmationLink": app.config.BaseUrl + "#/signup-confirm/" + token.plain,
		}

		err = app.mailer.Send(newUser.Email, "signup-confirm.tmpl", data)
		if err != nil {
			app.logger.Error(err)
		}
	})

	w.WriteHeader(http.StatusOK)
}

func (app *application) signupConfirmHandler(w http.ResponseWriter, r *http.Request) {
	token, err := app.readString(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userId, err := app.getAppUserIdFromToken(r.Context(), scopeSignup, token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if userId == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userId)).UpdateAll(r.Context(), app.db,
		models.M{models.AppUserColumns.Activated: true, models.AppUserColumns.LastAccess: time.Now()})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.deleteAllTokensForUser(r.Context(), userId, scopeSignup)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
