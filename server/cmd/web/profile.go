package main

import (
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"net/http"
)

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
