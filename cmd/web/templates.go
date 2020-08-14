package main

import (
	"dvhthomas/snippetbox/pkg/forms"
	"dvhthomas/snippetbox/pkg/models"
	"html/template"
	"net/url"
	"path/filepath"
	"time"
)

type templateData struct {
	Flash           string
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	FormData        url.Values
	Form            *forms.Form
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	// This is a weird Go-ism where Time.Format lets you define
	// a string using a specific agreed date (this one!)
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// We're using a longer contstuctor so that we can also pass in the map
		// of named functions defined above.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// This part really confused me until I realized that `ts` here is not
		// new, but is being redefined by adding more templates to the existing
		// to the *Template. The '=' instead of ':=' finally helped me get it!
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
