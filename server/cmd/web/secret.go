package main

import (
	"fmt"
	"net/http"
)

func (app *application) secret(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}
