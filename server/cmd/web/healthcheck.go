package main

import (
	"net/http"

	"gotodo.rasc.ch/cmd/web/output"
	"gotodo.rasc.ch/internal/response"
	"gotodo.rasc.ch/internal/version"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, output.HealthcheckOutput{
		Status: "up",
	})
}

func (app *application) appVersionHandler(w http.ResponseWriter, _ *http.Request) {
	v := version.Get()
	response.JSON(w, http.StatusOK, output.AppVersionOutput{
		BuildTime: v.BuildTime,
		Version:   v.Version,
	})
}
