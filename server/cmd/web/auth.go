package main

import (
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"log"
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
		log.Fatal(err)
	}
}

func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	value := app.sessionManager.Get(r.Context(), "userId")
	userId, ok := value.(int64)
	if !ok {
		userId = 0
	}

	if userId > 0 {
		ctx, cancel := app.createDbContext()
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
			app.writeJSON(w, r, http.StatusOK, map[string]interface{}{
				"authority": user.Authority,
			})
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

	var loginInput struct {
		Password string `name:"password" validate:"required,gte=8"`
		Username string `name:"username" validate:"required,email"`
	}
	err = app.decoder.Decode(&loginInput, r.PostForm)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	valid, fieldErrors := app.validate(loginInput)
	if !valid {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, UpdateResponse{
			Success:     false,
			FieldErrors: fieldErrors,
		})
		return
	}

	ctx, cancel := app.createDbContext()
	defer cancel()
	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Authority,
		models.AppUserColumns.PasswordHash,
		models.AppUserColumns.Expired,
		models.AppUserColumns.Activated),
		models.AppUserWhere.Email.EQ(loginInput.Username)).One(ctx, app.db)
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
			ctxUpdate, cancelUpdate := app.createDbContext()
			defer cancelUpdate()

			err := models.AppUsers(models.AppUserWhere.ID.EQ(user.ID)).UpdateAll(ctxUpdate, app.db,
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
	}
	w.WriteHeader(http.StatusOK)
}
