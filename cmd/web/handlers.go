package main

import (
	"dvhthomas/snippetbox/pkg/forms"
	"dvhthomas/snippetbox/pkg/models"
	"errors"
	"fmt"

	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// The pat library automatically handles a trailing '/' on the path
	// so this handler covers http://website _and_ http://website/
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// If the flash data exists this will get the value and remove it,
	// or return an empty string.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}

// Create a snippet page with a form. This could have pre-existing form
// data if the page is displaying errors and prior data from a failed POST.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// This only handles POST requests - look in routes.go for details
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		// Argh! Forgot the return here and couldn't figure
		// out why errors were getting through. It's because the
		// rest of the DB insert logic was running regardless.
		// Watch out for missing 'returns'!
		return
	}

	id, err := app.snippets.Insert(
		form.Get("title"),
		form.Get("content"),
		form.Get("expires"),
	)

	if err != nil {
		app.serverError(w, err)
	}

	// If there's no existing session for the user, the middleware will create
	// the session cookie automatically and *then* put the data in there.
	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
