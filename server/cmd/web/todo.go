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

	if ok := app.parseFromJson(w, r, &todoInput); !ok {
		return
	}

	userId := app.sessionManager.Get(r.Context(), "userId").(int64)

	var response interface{}
	var httpStatus int
	var err error

	ctx, cancel := app.createDbContext()
	if todoInput.Id > 0 {
		err = models.Todos(models.TodoWhere.ID.EQ(todoInput.Id), models.TodoWhere.AppUserID.EQ(userId)).
			UpdateAll(ctx, app.db, models.M{models.TodoColumns.Subject: todoInput.Subject,
				models.TodoColumns.Description: app.newNullString(todoInput.Description)})
		httpStatus = http.StatusOK
	} else {
		newTodo := models.Todo{
			Subject:     todoInput.Subject,
			Description: app.newNullString(todoInput.Description),
			AppUserID:   userId,
		}
		err = newTodo.Insert(ctx, app.db, boil.Infer())
		response = newTodo.ID
		httpStatus = http.StatusCreated
	}
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if response != nil {
		app.writeJSON(w, r, httpStatus, response)
	} else {
		w.WriteHeader(httpStatus)
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
	err = models.Todos(models.TodoWhere.ID.EQ(int64(todoId))).DeleteAll(ctx, app.db)
	cancel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
