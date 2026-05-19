package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.abdulmoiz.net/internal/models"
)



func (app *application)home(w http.ResponseWriter, r *http.Request) {
	
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return 
	}
	data := app.newTemplateData(r)
	data.Snippets=snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }
	// files := []string{
	// "./ui/html/base.tmpl",
	// "./ui/html/partials/nav.tmpl",
	// "./ui/html/pages/home.tmpl",
	// }
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// data := &templateData{
	// 	Snippets: snippets,
	// }
	// err = ts.ExecuteTemplate(w, "base",data)
	// if err !=nil {
	// 	app.serverError(w, err)

	// }
	
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	// app.infoLog.Printf("%d", id)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }
	// ts ,err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w,err)
	// 	return
	// }
	// data := &templateData{Snippet: snippet}
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// app.serverError(w, err)
	// }

}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}
func (app *application)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"	
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// w.Write([]byte("Create a new snippet..."))
}