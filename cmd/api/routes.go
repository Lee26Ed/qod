package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {

   // setup a new router
   router := httprouter.New()
   // handle 404
   router.NotFound = http.HandlerFunc(a.notFoundResponse)
  // handle 405
   router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
   // setup routes
   router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
   router.HandlerFunc(http.MethodPost, "/v1/quotes", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quotes/:id", a.displayQuoteHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/quotes/:id", a.updateQuoteHandler)

   return a.recoverPanic(router) 
}
