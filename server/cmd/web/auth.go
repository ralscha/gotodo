package main

import (
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/gobuffalo/validate"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/cmd/web/input"
	"gotodo.rasc.ch/cmd/web/output"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/models"
	"gotodo.rasc.ch/internal/request"
	"gotodo.rasc.ch/internal/response"
	"log/slog"
	"net/http"
	"time"
)

var userNotFoundPasswordHash string

func initAuth(config config.Config) error {
	var err error
	userNotFoundPasswordHash, err = argon2id.CreateHash("userNotFoundPassword", &argon2id.Params{
		Memory:      config.Argon2.Memory,
		Iterations:  config.Argon2.Iterations,
		Parallelism: config.Argon2.Parallelism,
		SaltLength:  config.Argon2.SaltLength,
		KeyLength:   config.Argon2.KeyLength,
	})
	return err
}

func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	if userID > 0 {
		user, err := models.AppUsers(qm.Select(
			models.AppUserColumns.Authority,
			models.AppUserColumns.Expired,
			models.AppUserColumns.Activated),
			models.AppUserWhere.ID.EQ(userID)).One(r.Context(), app.db)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			response.InternalServerError(w, err)
			return
		}
		if user != nil && user.Activated && user.Expired.IsZero() {
			response.JSON(w, http.StatusOK, output.LoginOutput{Authority: user.Authority})
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		response.InternalServerError(w, err)
	}

	var loginInput input.LoginInput
	if ok := request.DecodeJSONValidate(w, r, &loginInput); !ok {
		return
	}

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Authority,
		models.AppUserColumns.PasswordHash,
		models.AppUserColumns.Expired,
		models.AppUserColumns.Activated),
		models.AppUserWhere.Email.EQ(loginInput.Email)).One(r.Context(), app.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		response.InternalServerError(w, err)
		return
	}

	if user != nil && user.Activated && user.Expired.IsZero() {
		match, err := argon2id.ComparePasswordAndHash(loginInput.Password, user.PasswordHash)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		if match {
			err := models.AppUsers(models.AppUserWhere.ID.EQ(user.ID)).UpdateAll(r.Context(), app.db,
				models.M{models.AppUserColumns.LastAccess: time.Now()})
			if err != nil {
				response.InternalServerError(w, err)
				return
			}

			app.sessionManager.Put(r.Context(), "userID", user.ID)

			response.JSON(w, http.StatusOK, output.LoginOutput{Authority: user.Authority})
			return
		}
	} else {
		if _, err := argon2id.ComparePasswordAndHash(loginInput.Password, userNotFoundPasswordHash); err != nil {
			response.InternalServerError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.Destroy(r.Context()); err != nil {
		response.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) passwordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	var passwordResetRequestInput input.PasswordResetRequestInput
	if ok := request.DecodeJSONValidate(w, r, &passwordResetRequestInput); !ok {
		return
	}

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID),
		models.AppUserWhere.Email.EQ(passwordResetRequestInput.Email)).One(r.Context(), app.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		response.InternalServerError(w, err)
		return
	}

	if user != nil {
		token, err := app.insertToken(r.Context(), user.ID, app.config.Cleanup.PasswordResetTokenMaxAge, models.TokensScopePasswordReset)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		app.background(func() {
			data := map[string]any{
				"resetLink": app.config.BaseURL + "#/password-reset/" + token.plain,
			}

			err = app.mailer.Send(passwordResetRequestInput.Email, "password-reset.tmpl", data)
			if err != nil {
				slog.Error("sending password reset email failed", err)
			}
		})
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) passwordResetHandler(w http.ResponseWriter, r *http.Request) {
	var passwordResetInput input.PasswordResetInput
	if ok := request.DecodeJSONValidate(w, r, &passwordResetInput); !ok {
		return
	}

	userID, err := app.getAppUserIDFromToken(r.Context(), models.TokensScopePasswordReset, passwordResetInput.ResetToken)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if userID == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	compromised, err := app.isPasswordCompromised(passwordResetInput.Password)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if compromised {
		validationError := validate.Errors{
			Errors: map[string][]string{"password": {"weak"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	newPasswordHash, err := argon2id.CreateHash(passwordResetInput.Password, &argon2id.Params{
		Memory:      app.config.Argon2.Memory,
		Iterations:  app.config.Argon2.Iterations,
		Parallelism: app.config.Argon2.Parallelism,
		SaltLength:  app.config.Argon2.SaltLength,
		KeyLength:   app.config.Argon2.KeyLength,
	})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).UpdateAll(r.Context(), app.db,
		models.M{models.AppUserColumns.PasswordHash: newPasswordHash})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = app.deleteAllTokensForUser(r.Context(), userID, models.TokensScopePasswordReset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

}
