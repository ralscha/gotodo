package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"gotodo.rasc.ch/internal/config"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)

	if app.config.Environment == config.Development {
		router.Use(middleware.Logger)
	}

	router.Use(httprate.LimitAll(1_000, 1*time.Minute))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Minute))
	router.Use(middleware.NoCache)

	router.Route("/v1", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Get("/healthcheck", app.healthcheckHandler)
		r.Post("/authenticate", app.authenticateHandler)
		r.Post("/login", app.loginHandler)
		r.Post("/logout", app.logoutHandler)
		r.Mount("/", app.authenticatedRouter())
	})

	return router
}

func (app *application) authenticatedRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(app.authenticatedOnly)
	r.Get("/todo", app.todoGetHandler)
	r.Post("/todo", app.todoSaveHandler)
	r.Delete("/todo/{todoId:\\d+}", app.todoDeleteHandler)
	return r
}

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := app.sessionManager.Get(r.Context(), "userId")
		userId, ok := value.(int64)
		if ok && userId > 0 {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}
