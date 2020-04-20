package main

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/home", homeHandler).Methods("GET")
	r.HandleFunc("/blog", blogHandler).Methods("GET")
	r.HandleFunc("/contact", contactHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}
