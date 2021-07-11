package main

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"gotodo.rasc.ch/internal/middleware"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(chimiddleware.RealIP)
	router.Use(middleware.LoggerMiddleware())
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(time.Minute))

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
	r.Get("/secret", app.secret)
	return r
}

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := app.sessionManager.Get(r.Context(), "userId")
		userId, ok := value.(int64)
		if !ok {
			userId = 0
		}
		if userId > 0 {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}
