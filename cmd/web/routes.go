package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"snippetbox.abdulmoiz.net/ui"

	"github.com/justinas/alice"
)

func(app *application)routes() http.Handler {
	router := httprouter.New()
	// basically we are creating handler which wraps our notFound func so any 404 error will now call notFound()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	app.notFound(w)
	})
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.HandlerFunc(http.MethodGet, "/ping", ping)
	
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create",
	protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.accountView))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
	// first log Request middleware is exuceted then secure headers and finally main mux handler

	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

}