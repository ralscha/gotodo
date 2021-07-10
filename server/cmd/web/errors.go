package main

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	log.Error().Err(err).Str("request_method", r.Method).Str("request_url", r.URL.String()).Msg("")
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := anyMap{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
