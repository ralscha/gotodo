package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/models"
	"gotodo.rasc.ch/internal/response"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(response.NotFound)
	mux.MethodNotAllowed(response.MethodNotAllowed)

	// Middleware
	mux.Use(middleware.ClientIPFromXFF())
	if app.config.Environment == config.Development {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)
	mux.Use(httprate.LimitBy(1_000, 1*time.Minute, clientIPRateLimitKey))
	mux.Use(middleware.Timeout(15 * time.Second))
	mux.Use(middleware.NoCache)

	mux.Route("/v1", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Get("/healthcheck", app.healthcheckHandler)
		r.Post("/authenticate", app.authenticateHandler)
		r.Post("/login", app.loginHandler)
		r.Post("/signup", app.signupHandler)
		r.Post("/signup-confirm", app.signupConfirmHandler)
		r.Post("/password-reset-request", app.passwordResetRequestHandler)
		r.Post("/password-reset", app.passwordResetHandler)
		r.Mount("/", app.authenticatedRouter())
	})

	return mux
}

func clientIPRateLimitKey(r *http.Request) (string, error) {
	return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
}

func (app *application) authenticatedRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(app.authenticatedOnly)
	r.Post("/logout", app.logoutHandler)
	r.Get("/todo", app.todoGetHandler)
	r.Post("/todo", app.todoSaveHandler)
	r.Delete("/todo/{todoID:\\d+}", app.todoDeleteHandler)
	r.Get("/profile/build-info", app.appVersionHandler)
	r.Post("/profile/email-change", app.emailChangeHandler)
	r.Post("/profile/email-change-confirm", app.emailChangeConfirmHandler)
	r.Post("/profile/password-change", app.passwordChangeHandler)
	r.Post("/profile/account-delete", app.accountDeleteHandler)
	return r
}

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := app.sessionManager.GetInt64(r.Context(), "userID")
		if userID <= 0 {
			response.Unauthorized(w)
			return
		}

		user, err := models.AppUsers(
			qm.Select(
				models.AppUserColumns.PasswordHash,
				models.AppUserColumns.Activated,
				models.AppUserColumns.Expired,
			),
			models.AppUserWhere.ID.EQ(userID),
		).One(r.Context(), app.db)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			response.InternalServerError(w, err)
			return
		}
		sessionPasswordHash := app.sessionManager.GetString(r.Context(), "passwordHash")
		if user == nil || !user.Activated || !user.Expired.IsZero() ||
			sessionPasswordHash == "" || sessionPasswordHash != user.PasswordHash {
			if err := app.sessionManager.Destroy(r.Context()); err != nil {
				response.InternalServerError(w, err)
				return
			}
			response.Unauthorized(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
