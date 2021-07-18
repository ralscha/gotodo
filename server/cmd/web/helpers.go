package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null/v8"
	"io"
	"net/http"
	"strings"
	"time"
)

func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

const maxBytes = 1_048_576

func (app *application) readString(w http.ResponseWriter, r *http.Request) (string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) newNullString(s string) null.String {
	if len(s) == 0 {
		return null.String{}
	}
	return null.NewString(s, true)
}

func (app *application) validate(obj interface{}) (bool, map[string]string) {
	err := app.validator.Struct(obj)
	if err != nil {
		fieldErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			fieldErrors[err.Field()] = err.Tag()
		}
		return false, fieldErrors
	}
	return true, nil
}

func (app *application) createDbContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
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
