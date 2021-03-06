package main

import (
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/models"
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
	value := app.sessionManager.Get(r.Context(), "userId")
	userId, ok := value.(int64)
	if !ok {
		userId = 0
	}

	if userId > 0 {
		user, err := models.AppUsers(qm.Select(
			models.AppUserColumns.Authority,
			models.AppUserColumns.Expired,
			models.AppUserColumns.Activated),
			models.AppUserWhere.ID.EQ(userId)).One(r.Context(), app.db)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			app.serverErrorResponse(w, r, err)
			return
		}
		if user != nil && user.Activated && user.Expired.IsZero() {
			app.writeJSON(w, r, http.StatusOK, map[string]interface{}{
				"authority": user.Authority,
			})
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	var loginInput struct {
		Password string `name:"password" validate:"required,gte=8"`
		Email    string `name:"email" validate:"required,email"`
	}

	if ok := app.parseFromForm(w, r, &loginInput); !ok {
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
		app.serverErrorResponse(w, r, err)
		return
	}

	if user != nil && user.Activated && user.Expired.IsZero() {
		match, err := argon2id.ComparePasswordAndHash(loginInput.Password, user.PasswordHash)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if match {
			err := models.AppUsers(models.AppUserWhere.ID.EQ(user.ID)).UpdateAll(r.Context(), app.db,
				models.M{models.AppUserColumns.LastAccess: time.Now()})
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			app.sessionManager.Put(r.Context(), "userId", user.ID)

			app.writeJSON(w, r, http.StatusOK, map[string]interface{}{
				"authority": user.Authority,
			})
			return
		}
	} else {
		_, err := argon2id.ComparePasswordAndHash(loginInput.Password, userNotFoundPasswordHash)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.Destroy(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
