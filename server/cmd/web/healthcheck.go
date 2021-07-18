package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, r, http.StatusOK, map[string]interface{}{
		"status": "up",
	})
}

func (app *application) buildInfoHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, r, http.StatusOK, map[string]interface{}{
		"buildTime": appBuildTime,
		"version":   appVersion,
	})
}
