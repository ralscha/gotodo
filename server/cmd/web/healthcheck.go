package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "available",
		"system_info": map[string]string{
			"buildTime": appBuildTime,
			"version":   appVersion,
		},
	}

	app.writeJSON(w, r, http.StatusOK, response)
}
