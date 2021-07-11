package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"net/http"
	"time"
)

var userNotFoundPasswordHash string

func init() {
	params := &argon2id.Params{
		Memory:      1 << 17,
		Iterations:  20,
		Parallelism: 8,
		SaltLength:  16,
		KeyLength:   32,
	}
	var err error
	userNotFoundPasswordHash, err = argon2id.CreateHash("userNotFoundPassword", params)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	value := app.sessionManager.Get(r.Context(), "userId")
	userId, ok := value.(int64)
	if !ok {
		userId = 0
	}

	if userId > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		user, err := models.AppUsers(qm.Select(
			models.AppUserColumns.Authority,
			models.AppUserColumns.Expired,
			models.AppUserColumns.Activated),
			models.AppUserWhere.ID.EQ(userId)).One(ctx, app.db)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			app.serverErrorResponse(w, r, err)
			return
		}
		if user != nil && user.Activated && user.Expired.IsZero() {
			data := anyMap{
				"authority": user.Authority,
			}
			err := app.writeJSON(w, http.StatusOK, data, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}

		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	inputEmail := r.PostForm.Get("username")
	inputPassword := r.PostForm.Get("password")
	if inputEmail == "" || inputPassword == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Authority,
		models.AppUserColumns.PasswordHash,
		models.AppUserColumns.Expired,
		models.AppUserColumns.Activated),
		models.AppUserWhere.Email.EQ(inputEmail)).One(ctx, app.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.serverErrorResponse(w, r, err)
		return
	}

	if user != nil && user.Activated && user.Expired.IsZero() {
		match, err := argon2id.ComparePasswordAndHash(inputPassword, user.PasswordHash)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if match {
			app.sessionManager.Put(r.Context(), "userId", user.ID)

			data := anyMap{
				"authority": user.Authority,
			}
			err := app.writeJSON(w, http.StatusOK, data, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		}
	} else {
		_, err := argon2id.ComparePasswordAndHash(inputPassword, userNotFoundPasswordHash)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.Destroy(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}