package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotodo.rasc.ch/internal/models"
	"net/http"
	"strconv"
)

func (app *application) todoGetHandler(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	ctx, cancel := app.createDbContext()
	todos, err := models.Todos(models.TodoWhere.AppUserID.EQ(userId)).All(ctx, app.db)
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if todos != nil {
		app.writeJSON(w, r, http.StatusOK, todos)
	} else {
		app.writeJSON(w, r, http.StatusOK, make([]interface{}, 0))
	}
}

func (app *application) todoInsertHandler(w http.ResponseWriter, r *http.Request) {
	var todoInput struct {
		Subject     string `name:"subject" validate:"required"`
		Description string
	}

	err := app.readJSON(w, r, &todoInput)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	valid, fieldErrors := app.validate(todoInput)
	if !valid {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, InsertResponse{
			UpdateResponse: UpdateResponse{
				Success:     false,
				FieldErrors: fieldErrors,
			},
		})
		return
	}

	newTodo := models.Todo{
		Subject:     todoInput.Subject,
		Description: app.newNullString(todoInput.Description),
		AppUserID:   app.sessionManager.Get(r.Context(), "userId").(int64),
	}

	ctx, cancel := app.createDbContext()
	err = newTodo.Insert(ctx, app.db, boil.Infer())
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, r, http.StatusOK, InsertResponse{
		Id: newTodo.ID,
		UpdateResponse: UpdateResponse{
			Success: true,
		},
	})
}

func (app *application) todoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var todoInput struct {
		Id          int64  `name:"id" validate:"gt=0"`
		Subject     string `name:"subject" validate:"required"`
		Description string
	}
	err := app.readJSON(w, r, &todoInput)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	valid, fieldErrors := app.validate(todoInput)
	if !valid {
		app.writeJSON(w, r, http.StatusUnprocessableEntity, UpdateResponse{
			Success:     false,
			FieldErrors: fieldErrors,
		})
		return
	}

	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	ctx, cancel := app.createDbContext()
	err = models.Todos(models.TodoWhere.ID.EQ(todoInput.Id), models.TodoWhere.AppUserID.EQ(userId)).
		UpdateAll(ctx, app.db, models.M{models.TodoColumns.Subject: todoInput.Subject,
			models.TodoColumns.Description: app.newNullString(todoInput.Description)})
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, r, http.StatusOK, UpdateResponse{
		Success: true,
	})
}

func (app *application) todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	todoIdStr := chi.URLParam(r, "todoId")
	todoId, err := strconv.Atoi(todoIdStr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ctx, cancel := app.createDbContext()
	err = models.Todos(models.TodoWhere.ID.EQ(int64(todoId))).DeleteAll(ctx, app.db)
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	deleteResponse := DeleteResponse{
		Success: true,
	}
	app.writeJSON(w, r, http.StatusOK, deleteResponse)
}
