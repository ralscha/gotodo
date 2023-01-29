package main

import (
	"github.com/alexedwards/argon2id"
	"github.com/gobuffalo/validate"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/exp/slog"
	"gotodo.rasc.ch/cmd/web/input"
	"gotodo.rasc.ch/internal/models"
	"gotodo.rasc.ch/internal/request"
	"gotodo.rasc.ch/internal/response"
	"net/http"
	"time"
)

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var signUpInput input.SignUpInput
	if ok := request.DecodeJSONValidate(w, r, &signUpInput); !ok {
		return
	}

	exists, err := models.AppUsers(
		models.AppUserWhere.Email.EQ(signUpInput.Email),
		qm.Or2(models.AppUserWhere.EmailNew.EQ(null.NewString(signUpInput.Email, true))),
	).Exists(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if exists {
		validationError := validate.Errors{
			Errors: map[string][]string{"email": {"exists"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	compromised, err := app.isPasswordCompromised(signUpInput.Password)
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

	hash, err := argon2id.CreateHash(signUpInput.Password, &argon2id.Params{
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

	newUser := models.AppUser{
		Email:        signUpInput.Email,
		PasswordHash: hash,
		Authority:    "USER",
		Activated:    false,
	}

	err = newUser.Insert(r.Context(), app.db, boil.Infer())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	token, err := app.insertToken(r.Context(), newUser.ID, app.config.Cleanup.SignupTokenMaxAge, models.TokensScopeSignup)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"confirmationLink": app.config.BaseURL + "#/signup-confirm/" + token.plain,
		}

		err = app.mailer.Send(newUser.Email, "signup-confirm.tmpl", data)
		if err != nil {
			slog.Error("sending signup confirmation email failed", err)
		}
	})

	w.WriteHeader(http.StatusOK)
}

func (app *application) signupConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var tokenInput input.TokenInput
	if ok := request.DecodeJSONValidate(w, r, &tokenInput); !ok {
		return
	}

	userID, err := app.getAppUserIDFromToken(r.Context(), models.TokensScopeSignup, tokenInput.Token)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if userID == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).UpdateAll(r.Context(), app.db,
		models.M{models.AppUserColumns.Activated: true, models.AppUserColumns.LastAccess: time.Now()})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = app.deleteAllTokensForUser(r.Context(), userID, models.TokensScopeSignup)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
