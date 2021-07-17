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

func (app *application) todoSaveHandler(w http.ResponseWriter, r *http.Request) {
	var todoInput struct {
		Id          int64
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
		app.writeJSON(w, r, http.StatusUnprocessableEntity, SaveResponse{
			Success:     false,
			FieldErrors: fieldErrors,
		})
		return
	}

	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	var saveResponse SaveResponse
	ctx, cancel := app.createDbContext()
	if todoInput.Id > 0 {
		err = models.Todos(models.TodoWhere.ID.EQ(todoInput.Id), models.TodoWhere.AppUserID.EQ(userId)).
			UpdateAll(ctx, app.db, models.M{models.TodoColumns.Subject: todoInput.Subject,
				models.TodoColumns.Description: app.newNullString(todoInput.Description)})
		saveResponse.Success = true
	} else {
		newTodo := models.Todo{
			Subject:     todoInput.Subject,
			Description: app.newNullString(todoInput.Description),
			AppUserID:   userId,
		}
		err = newTodo.Insert(ctx, app.db, boil.Infer())
		saveResponse.Id = newTodo.ID
		saveResponse.Success = true
	}
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, r, http.StatusOK, saveResponse)
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
