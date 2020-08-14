package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// Without the call depth of 2, the default log message will
	// always show that helpers.go is the source of the error, whereas
	// we want one level back from the helper file.
	app.errorLog.Output(2, trace)

	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the provided
	// name, call the serverError helper.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		// Forgot the return statement previously, so I got the 500 error as expected,
		// then baffled why the bad HTML content/error still rendered. Obviously
		// it's because I also called buf.WriteTo(w) as well even in the bad case. Doh!
		return
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	// Add the flash message to the template data if one exists.
	// Using PopString will ensure that it's a one-time use thing.
	td.Flash = app.session.PopString(r, "flash")

	// Add authenticated status to the template data
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

// Return true if the current *request* is from an authenticated user, otherwise
// return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	// The type assertion to bool here will default to false if
	// no value is found. And we return false for error cases.
	// So basically we're playing it safe with false unless we're positive
	// it's true: we assume there IS NOT an authenticated user.
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
