package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	// Setup routes and corresponding handlers
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /shigoto/view/{id}", app.shigotoView)
	mux.HandleFunc("GET /shigoto/create", app.shigotoCreate)
	mux.HandleFunc("POST /shigoto/create", app.shigotoCreatePost)
	return mux
}
