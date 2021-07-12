package main

import (
	"encoding/json"
	"net/http"
)

type anyMap map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data anyMap, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {

		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error("background recover", err)
			}
		}()

		fn()
	}()
}
