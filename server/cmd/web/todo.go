package main

import (
	"net/http"
	"strconv"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/go-chi/chi/v5"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"gotodo.rasc.ch/internal/models"
	"gotodo.rasc.ch/internal/request"
	"gotodo.rasc.ch/internal/response"
)

type ValidatedTodo models.Todo

func (v *ValidatedTodo) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "subject",
			Field:   v.Subject,
			Message: "required",
		},
		&validators.StringLengthInRange{
			Name:    "subject",
			Field:   v.Subject,
			Message: "lte",
			Min:     0,
			Max:     255,
		},
		&validators.StringLengthInRange{
			Name:    "description",
			Field:   v.Description.String,
			Message: "lte",
			Min:     0,
			Max:     255,
		},
	)
}

func (app *application) todoGetHandler(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	todos, err := models.Todos(
		models.TodoWhere.AppUserID.EQ(userID),
		qm.OrderBy(models.TodoColumns.ID+" ASC"),
	).All(r.Context(), app.db)
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

	if todoInput.ID > 0 {
		result, err := app.db.ExecContext(r.Context(), `
			UPDATE todo
			SET subject = $1, description = $2
			WHERE id = $3 AND app_user_id = $4`,
			todoInput.Subject, todoInput.Description, todoInput.ID, userID)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		if rowsAffected == 0 {
			response.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	newTodo := models.Todo{
		Subject:     todoInput.Subject,
		Description: todoInput.Description,
		AppUserID:   userID,
	}
	if err := newTodo.Insert(r.Context(), app.db, boil.Infer()); err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, models.Todo{ID: newTodo.ID})
}

func (app *application) todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	todoIDStr := chi.URLParam(r, "todoID")
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.NotFound(w, r)
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	result, err := app.db.ExecContext(r.Context(),
		"DELETE FROM todo WHERE id = $1 AND app_user_id = $2", todoID, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if rowsAffected == 0 {
		response.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
