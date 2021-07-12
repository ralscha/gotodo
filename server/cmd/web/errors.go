package main

import (
	"go.uber.org/zap"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Errorw("", zap.Error(err), zap.String("request_method", r.Method),
		zap.String("request_url", r.URL.String()))
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	w.WriteHeader(http.StatusInternalServerError)
}
