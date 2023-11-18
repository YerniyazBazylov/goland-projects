package main

import (
	"assignment2.yerniyaz.net/internal/data"
	"assignment2.yerniyaz.net/internal/validator"
	"fmt"
	"net/http"
	"time"
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
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showRemoteCarsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	RemoteCars := data.RemoteCars{
		ID:          id,
		CreatedAt:   time.Now(),
		Name:        "Remote Car BMW extra luxe",
		Year:        1989,
		Cost:        53,
		Description: "Remote Car made in Germany. ",
		Version:     1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"RemoteCars": RemoteCars}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
