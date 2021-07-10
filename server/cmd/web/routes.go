package main

import (
	"fmt"
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

	router.Get("/", greet)
	router.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)
	})

	return router
}

func greet(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}
