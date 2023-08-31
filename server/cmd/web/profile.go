package main

import (
	"github.com/alexedwards/argon2id"
	"github.com/gobuffalo/validate"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/cmd/web/input"
	"gotodo.rasc.ch/internal/models"
	"gotodo.rasc.ch/internal/request"
	"gotodo.rasc.ch/internal/response"
	"log/slog"
	"net/http"
)

func (app *application) emailChangeHandler(w http.ResponseWriter, r *http.Request) {
	var emailChangeInput input.EmailChangeInput
	if ok := request.DecodeJSONValidate(w, r, &emailChangeInput); !ok {
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.PasswordHash),
		models.AppUserWhere.ID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(emailChangeInput.Password, user.PasswordHash)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if !match {
		validationError := validate.Errors{
			Errors: map[string][]string{"password": {"invalid"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	exists, err := models.AppUsers(
		models.AppUserWhere.Email.EQ(emailChangeInput.NewEmail),
		qm.Or2(models.AppUserWhere.EmailNew.EQ(null.NewString(emailChangeInput.NewEmail, true))),
	).Exists(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if exists {
		validationError := validate.Errors{
			Errors: map[string][]string{"newEmail": {"exists"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).UpdateAll(r.Context(), app.db,
		models.M{models.AppUserColumns.EmailNew: emailChangeInput.NewEmail})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	token, err := app.insertToken(r.Context(), userID, app.config.Cleanup.EmailChangeTokenMaxAge, models.TokensScopeEmailChange)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"confirmationLink": app.config.BaseURL + "#/profile/email-confirm/" + token.plain,
		}

		err = app.mailer.Send(emailChangeInput.NewEmail, "email-change.tmpl", data)
		if err != nil {
			slog.Error("sending email confirm email failed", err)
		}
	})

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) emailChangeConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var tokenInput input.TokenInput
	if ok := request.DecodeJSONValidate(w, r, &tokenInput); !ok {
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	userIDFromToken, err := app.getAppUserIDFromToken(r.Context(), models.TokensScopeEmailChange, tokenInput.Token)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if userID != userIDFromToken {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.EmailNew),
		models.AppUserWhere.ID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if !user.EmailNew.Valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).UpdateAll(r.Context(), app.db,
		models.M{models.AppUserColumns.Email: user.EmailNew,
			models.AppUserColumns.EmailNew: null.NewString("", false)})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = app.deleteAllTokensForUser(r.Context(), userID, models.TokensScopeEmailChange)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) passwordChangeHandler(w http.ResponseWriter, r *http.Request) {
	var passwordChangeInput input.PasswordChangeInput
	if ok := request.DecodeJSONValidate(w, r, &passwordChangeInput); !ok {
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.PasswordHash),
		models.AppUserWhere.ID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(passwordChangeInput.OldPassword, user.PasswordHash)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if !match {
		validationError := validate.Errors{
			Errors: map[string][]string{"oldPassword": {"invalid"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	compromised, err := app.isPasswordCompromised(passwordChangeInput.NewPassword)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if compromised {
		validationError := validate.Errors{
			Errors: map[string][]string{"newPassword": {"weak"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	newPasswordHash, err := argon2id.CreateHash(passwordChangeInput.NewPassword, &argon2id.Params{
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

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) accountDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var passwordInput input.PasswordInput
	if ok := request.DecodeJSONValidate(w, r, &passwordInput); !ok {
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.PasswordHash),
		models.AppUserWhere.ID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(passwordInput.Password, user.PasswordHash)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if !match {
		validationError := validate.Errors{
			Errors: map[string][]string{"password": {"invalid"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	tx, err := app.db.BeginTx(r.Context(), nil)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.Todos(models.TodoWhere.AppUserID.EQ(userID)).DeleteAll(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.Tokens(models.TokenWhere.AppUserID.EQ(userID)).DeleteAll(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).DeleteAll(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = app.sessionManager.Destroy(r.Context())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
