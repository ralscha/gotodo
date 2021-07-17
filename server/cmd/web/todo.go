package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotodo.rasc.ch/internal/models"
	"net/http"
	"strconv"
)

func (app *application) todoGetHandler(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	ctx, cancel := app.createDbContext()
	defer cancel()
	todos, err := models.Todos(models.TodoWhere.AppUserID.EQ(userId)).All(ctx, app.db)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if todos != nil {
		err = app.writeJSON(w, http.StatusOK, todos, nil)
	} else {
		err = app.writeJSON(w, http.StatusOK, make([]models.Todo, 0), nil)
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) todoInsertHandler(w http.ResponseWriter, r *http.Request) {

	var newTodo struct {
		Subject     string `json:"subject" validate:"required"`
		Description string
	}

	err := app.readJSON(w, r, &newTodo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var insertResponse InsertResponse

	err = app.validator.Struct(newTodo)
	if err != nil {
		insertResponse.FieldErrors = make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			insertResponse.FieldErrors[err.Field()] = err.Tag()
		}
		err := app.writeJSON(w, http.StatusUnprocessableEntity, insertResponse, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var todo models.Todo
	todo.Subject = newTodo.Subject
	todo.Description = app.newNullString(newTodo.Description)
	todo.AppUserID = app.sessionManager.Get(r.Context(), "userId").(int64)

	ctx, cancel := app.createDbContext()
	defer cancel()
	err = todo.Insert(ctx, app.db, boil.Infer())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	insertResponse.Id = todo.ID
	insertResponse.Success = true

	err = app.writeJSON(w, http.StatusOK, insertResponse, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) todoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateTodo struct {
		Id          int64  `json:"id" validate:"gt=0"`
		Subject     string `json:"subject" validate:"required"`
		Description string
	}
	err := app.readJSON(w, r, &updateTodo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var updateResponse UpdateResponse

	err = app.validator.Struct(updateTodo)
	if err != nil {
		updateResponse.FieldErrors = make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			updateResponse.FieldErrors[err.Field()] = err.Tag()
		}
		err := app.writeJSON(w, http.StatusUnprocessableEntity, updateResponse, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	ctx, cancel := app.createDbContext()
	defer cancel()
	err = models.Todos(models.TodoWhere.ID.EQ(updateTodo.Id), models.TodoWhere.AppUserID.EQ(userId)).
		UpdateAll(ctx, app.db, models.M{models.TodoColumns.Subject: updateTodo.Subject,
			models.TodoColumns.Description: app.newNullString(updateTodo.Description)})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	updateResponse.Success = true

	err = app.writeJSON(w, http.StatusOK, updateResponse, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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

	deleteResponse := DeleteResponse{
		Success: true,
	}
	err = app.writeJSON(w, http.StatusOK, deleteResponse, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
