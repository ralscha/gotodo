package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null/v8"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (app *application) isPasswordCompromised(password string) (bool, error) {
	alg := sha1.New()
	alg.Write([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(alg.Sum(nil)))
	prefix := strings.ToUpper(hash[:5])
	suffix := strings.ToUpper(hash[5:])

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.pwnedpasswords.com/range/"+prefix, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Add-Padding", "true")

	response, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected http status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	hashList := strings.Split(string(body), "\r\n")

	for _, line := range hashList {
		if line[:35] == suffix {
			occurrence, err := strconv.ParseInt(line[36:], 10, 64)
			if err != nil {
				return false, err
			}
			if occurrence > 0 {
				return true, nil
			}
			break
		}
	}

	return false, nil
}

func (app *application) parseFromForm(w http.ResponseWriter, r *http.Request, input interface{}) bool {
	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return false
	}

	err = app.decoder.Decode(input, r.PostForm)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return false
	}

	valid, fieldErrors := app.validate(input)
	if !valid {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: fieldErrors,
		})
		return false
	}

	return true
}

func (app *application) parseFromJson(w http.ResponseWriter, r *http.Request, input interface{}) bool {
	err := app.readJSON(w, r, input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return false
	}

	valid, fieldErrors := app.validate(input)
	if !valid {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, FormErrorResponse{
			FieldErrors: fieldErrors,
		})
		return false
	}

	return true
}

func (app *application) schedule(fn func(), delay time.Duration) chan struct{} {
	stop := make(chan struct{})

	go func() {
		for {
			fn()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
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
