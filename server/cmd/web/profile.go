package main

import (
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"net/http"
)

func (app *application) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OldPassword string `name:"oldPassword" validate:"required,gte=8"`
		NewPassword string `name:"newPassword" validate:"required,gte=8"`
	}
	if ok := app.parseFromForm(w, r, &input); !ok {
		return
	}

	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.PasswordHash),
		models.AppUserWhere.ID.EQ(userId)).One(r.Context(), app.db)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.OldPassword, user.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{"oldPassword": "invalid"},
		})
		return
	}

	compromised, err := app.isPasswordCompromised(input.NewPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if compromised {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{
				"newPassword": "weak",
			},
		})
		return
	}

	newPasswordHash, err := argon2id.CreateHash(input.NewPassword, &argon2id.Params{
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

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	password, err := app.readString(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if password == "" {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{"password": "required"},
		})
		return
	}

	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.PasswordHash),
		models.AppUserWhere.ID.EQ(userId)).One(r.Context(), app.db)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: map[string]string{"password": "invalid"},
		})
		return
	}

	tx, err := app.db.BeginTx(r.Context(), nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = models.Todos(models.TodoWhere.AppUserID.EQ(userId)).DeleteAll(r.Context(), tx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = models.Tokens(models.TokenWhere.AppUserID.EQ(userId)).DeleteAll(r.Context(), tx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userId)).DeleteAll(r.Context(), tx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.sessionManager.Destroy(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
