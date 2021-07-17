package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := map[string]interface{}{
		"status": "available",
		"system_info": map[string]string{
			"buildTime": appBuildTime,
			"version":   appVersion,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
