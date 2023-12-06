package main

import (
	"assignment3.yerniyaz.net/internal/data"
	"assignment3.yerniyaz.net/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createRemoteCarsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string    `json:"name"`
		Year        int32     `json:"year"`
		Cost        data.Cost `json:"cost"`
		Description string    `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	remotecars := &data.RemoteCars{
		Name:        input.Name,
		Year:        input.Year,
		Cost:        input.Cost,
		Description: input.Description,
	}

	v := validator.New()

	if data.ValidateRemoteCars(v, remotecars); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.RemoteCars.Insert(remotecars)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/remote-cars/%d", remotecars.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"remotecars": remotecars}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRemoteCarsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	remotecars, err := app.models.RemoteCars.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"classiccars": remotecars}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateRemoteCarsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	remotecars, err := app.models.RemoteCars.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string    `json:"name"`
		Year        *int32     `json:"year"`
		Cost        *data.Cost `json:"cost"`
		Description *string    `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		remotecars.Name = *input.Name
	}
	if input.Year != nil {
		remotecars.Year = *input.Year
	}
	if input.Cost != nil {
		remotecars.Cost = *input.Cost
	}
	if input.Description != nil {
		remotecars.Description = *input.Description
	}

	v := validator.New()
	if data.ValidateRemoteCars(v, remotecars); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.RemoteCars.Update(remotecars)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"remotecars": remotecars}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRemoteCarsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.RemoteCars.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "remote_car successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRemoteCarsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string
		Description string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "year", "cost", "-id", "-name", "-year", "-cost"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	remotecars, metadata, err := app.models.RemoteCars.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": remotecars, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
