package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gotodo.rasc.ch/internal/models"
	"net/http"
	"strconv"
)

func (app *application) todoGetHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}

func (app *application) todoInsertHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}

func (app *application) todoUpdateHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}

func (app *application) todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	todoIdStr := chi.URLParam(r, "todoId")
	todoId, err := strconv.Atoi(todoIdStr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ctx, cancel := app.createDbContext()
	defer cancel()

	err = models.Todos(models.TodoWhere.ID.EQ(int64(todoId))).DeleteAll(ctx, app.db)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
