package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
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
			data := map[string]interface{}{
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

type LoginForm struct {
	Password string `validate:"required,gte=8"`
	Username string `validate:"required,email"`
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

	var lf LoginForm
	err = app.decoder.Decode(&lf, r.PostForm)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.validator.Struct(lf)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		w.WriteHeader(http.StatusUnauthorized)
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
		models.AppUserWhere.Email.EQ(lf.Username)).One(ctx, app.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.serverErrorResponse(w, r, err)
		return
	}

	if user != nil && user.Activated && user.Expired.IsZero() {
		match, err := argon2id.ComparePasswordAndHash(lf.Password, user.PasswordHash)
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

			data := map[string]interface{}{
				"authority": user.Authority,
			}
			err = app.writeJSON(w, http.StatusOK, data, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	} else {
		_, err := argon2id.ComparePasswordAndHash(lf.Password, userNotFoundPasswordHash)
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
}
