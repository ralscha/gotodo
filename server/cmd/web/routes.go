package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/response"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(response.NotFound)
	mux.MethodNotAllowed(response.MethodNotAllowed)

	// Middleware
	mux.Use(middleware.RealIP)
	if app.config.Environment == config.Development {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)
	mux.Use(httprate.LimitAll(1_000, 1*time.Minute))
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
		value := app.sessionManager.Get(r.Context(), "userID")
		userID, ok := value.(int64)
		if ok && userID > 0 {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}
