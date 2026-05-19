package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/justinas/alice"
)

func(app *application)routes() http.Handler {
	router := httprouter.New()
	// basically we are creating handler which wraps our notFound func so any 404 error will now call notFound()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static",
	fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
	// first log Request middleware is exuceted then secure headers and finally main mux handler

	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

}