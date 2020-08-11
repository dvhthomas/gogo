package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	// This actually matches '/snippet/create' but would assign the value
	// 'create' to the id variable. Which isn't really what we want since there's
	// no snippet with the id 'create'. So put this _after_ the '/snippet/create'
	// pattern in our code.
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
