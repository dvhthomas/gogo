package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// All dynamic routes will have a session cookie courtesy of golangcollege,
	// and a CSRF cookie courtesy of noSurf. Then we add a context value to
	// show whether the user session includes an authenticated user.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	// We're adding the session middleware to all the routes...
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	// Add the authentication middleware. The Append method makes the requireAuthentication
	// call the final one in the chain (see https://godoc.org/github.com/justinas/alice#Chain.Append)
	mux.Get("/snippet/create", dynamicMiddleware.
		Append(app.requireAuthentication).
		ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.
		Append(app.requireAuthentication).
		ThenFunc(app.createSnippet))
	// This actually matches '/snippet/create' but would assign the value
	// 'create' to the id variable. Which isn't really what we want since there's
	// no snippet with the id 'create'. So put this _after_ the '/snippet/create'
	// pattern in our code.
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// User-related routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))

	// Don't want unauthenticated users getting to the logout page
	mux.Post("/user/logout", dynamicMiddleware.
		Append(app.requireAuthentication).
		ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	// ...but we're not adding the session middleware to static routes
	// because it's inherently stateless content. No cookie required!
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
