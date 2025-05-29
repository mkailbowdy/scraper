package main

import (
	"errors"
	"fmt"
	"html/template"
	"jobscraper.kailmendoza.com/internal/models"
	"net/http"
	"strconv"
)

// Define a handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	shigotos, err := app.shigotos.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Shigotos: shigotos,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) shigotoView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	shigoto, err := app.shigotos.Get(id)
	if err != nil {
		// Use our custom error models.ErrNoRecord. NOT sql.ErrorNoRows!
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// This templateData struct will hold all dynamic data we need for the html templates
	data := templateData{
		Shigoto: shigoto,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) shigotoCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create job"))
}

func (app *application) shigotoCreatePost(w http.ResponseWriter, r *http.Request) {
	companyName := "Samsung"
	jobTitle := "QA Engineer"
	category := "IT"
	location := "Seoul, Korea"
	employmentType := "Contract"
	description := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	japaneseLevel := "N1"
	englishLevel := "Native"
	sponsorship := true

	id, err := app.shigotos.Insert(companyName, jobTitle, category, location, employmentType, description, japaneseLevel, englishLevel, sponsorship)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/shigoto/view/%d", id), http.StatusSeeOther)
}
