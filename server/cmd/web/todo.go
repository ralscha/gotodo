package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotodo.rasc.ch/internal/models"
	"gotodo.rasc.ch/internal/request"
	"gotodo.rasc.ch/internal/response"
	"net/http"
	"strconv"
)

type ValidatedTodo models.Todo

func (v *ValidatedTodo) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "subject",
			Field:   v.Subject,
			Message: "required",
		},
	)
}

func (app *application) todoGetHandler(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	todos, err := models.Todos(models.TodoWhere.AppUserID.EQ(userID)).All(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if todos == nil {
		todos = []*models.Todo{}
	}
	response.JSON(w, http.StatusOK, todos)
}

func (app *application) todoSaveHandler(w http.ResponseWriter, r *http.Request) {
	var todoInput ValidatedTodo
	if ok := request.DecodeJSONValidate(w, r, &todoInput); !ok {
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	var newID int64
	var httpStatus int
	var err error

	if todoInput.ID > 0 {
		err = models.Todos(models.TodoWhere.ID.EQ(todoInput.ID), models.TodoWhere.AppUserID.EQ(userID)).
			UpdateAll(r.Context(), app.db, models.M{models.TodoColumns.Subject: todoInput.Subject,
				models.TodoColumns.Description: todoInput.Description})
		httpStatus = http.StatusOK
	} else {
		newTodo := models.Todo{
			Subject:     todoInput.Subject,
			Description: todoInput.Description,
			AppUserID:   userID,
		}
		err = newTodo.Insert(r.Context(), app.db, boil.Infer())
		newID = newTodo.ID
		httpStatus = http.StatusCreated
	}
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if newID > 0 {
		response.JSON(w, httpStatus, models.Todo{
			ID: newID,
		})
	} else {
		w.WriteHeader(httpStatus)
	}
}

func (app *application) todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	todoIDStr := chi.URLParam(r, "todoID")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.Todos(models.TodoWhere.ID.EQ(int64(todoID))).DeleteAll(r.Context(), app.db)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
