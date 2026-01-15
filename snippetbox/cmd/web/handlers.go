package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"snippetbox.alexedwards.net/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w,err)
		return 
	}

	data := app.newTemplateData(r)
	data.snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return 
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		}
		else{
			app.serverError(w,err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}


func (app *application) snippetCreatePost(w http.ResponseWriter, r
*http.Request) {
// Checking if the request method is a POST is now superfluous and can be
// removed, because this is done automatically by httprouter.
title := "O snail"
content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n–
Kobayashi Issa"
expires := 7
id, err := app.snippets.Insert(title, content, expires)
if err != nil {
app.serverError(w, err)
return
}
// Update the redirect path to use the new clean URL format.
http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id),
http.StatusSeeOther)
}