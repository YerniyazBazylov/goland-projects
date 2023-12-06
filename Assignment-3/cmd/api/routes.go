package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/remote-cars", app.listRemoteCarsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/remote-cars", app.createRemoteCarsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/remote-cars/:id", app.showRemoteCarsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/remote-cars/:id", app.updateRemoteCarsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/remote-cars/:id", app.deleteRemoteCarsHandler)

	return app.recoverPanic(app.rateLimit(router))

}
