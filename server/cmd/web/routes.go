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
	})

	return router
}
