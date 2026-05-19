package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.abdulmoiz.net/internal/models"
)
type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}



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
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}
func (app *application)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// title := "O snail"	
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := 7
	// id, err := app.snippets.Insert(title, content, expires)
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content:r.PostForm.Get("content"),
		Expires: expires,
		FieldErrors: map[string]string{},

	}
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// w.Write([]byte("Create a new snippet..."))
}