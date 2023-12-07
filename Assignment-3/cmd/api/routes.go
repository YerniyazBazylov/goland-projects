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

	router.HandlerFunc(http.MethodGet, "/v1/remote-cars", app.requirePermission("remote-cars:read", app.listRemoteCarsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/remote-cars", app.requirePermission("remote-cars:write", app.createRemoteCarsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/remote-cars/:id", app.requirePermission("remote-cars:read", app.showRemoteCarsHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/remote-cars/:id", app.requirePermission("remote-cars:write", app.updateRemoteCarsHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/remote-cars/:id", app.requirePermission("remote-cars:write", app.deleteRemoteCarsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}
