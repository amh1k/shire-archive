package main

import (
	"fmt"
	"net/http"
)

// It’s important to know that when the last handler in the chain returns, control is
// passed back up the chain in the reverse direction. So when our code is being
// executed the flow of control actually looks like this:
// secureHeaders → servemux → application handler → servemux → secureHeader
// In any middleware handler, code which comes before next.ServeHTTP() will be
// executed on the way down the chain, and any code after next.ServeHTTP() — or
// in a deferred function — will be executed on the way back up
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request, ) {
		w.Header().Set("Content-Security-Policy","default-src 'self'; style-src 'self' fonts.googleapis.com; font-srcfonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)

		//This code will execute on the way back up the chain
})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
		r.URL.RequestURI())
		next.ServeHTTP(w, r)

	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ){
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			if err := recover(); err !=  nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)


	})
}