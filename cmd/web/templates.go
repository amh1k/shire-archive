package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.abdulmoiz.net/internal/models"
	"snippetbox.abdulmoiz.net/ui"
)
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}
var functions= template.FuncMap{
	"humanDate" : humanDate,
}

func newTemplateCache() (map [string] *template.Template, error) {
	cache := map[string] *template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		
		
		cache[name] = ts
	}
	return cache, nil

}

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets [] *models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string

}
