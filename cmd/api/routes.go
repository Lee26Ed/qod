package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthcheckHandler)

	return router
}