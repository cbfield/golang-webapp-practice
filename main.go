package main

import (
	"net/http"
	"text/template"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("7075318008"))
var client *redis.Client
var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandlerGet).Methods("GET")
	r.HandleFunc("/", rootHandlerPost).Methods("POST")

	r.HandleFunc("/login", loginHandlerGet).Methods("GET")
	r.HandleFunc("/login", loginHandlerPost).Methods("POST")
	r.HandleFunc("/test", testHandler).Methods("GET")

	r.HandleFunc("/home", homeHandlerGet).Methods("GET")
	r.HandleFunc("/blog", blogHandlerGet).Methods("GET")
	r.HandleFunc("/contact", contactHandlerGet).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func rootHandlerGet(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}

func rootHandlerPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	client.LPush("comments", comment)
	http.Redirect(w, r, "/", 302)
}

func loginHandlerGet(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	untyped, err := session.Values["username"]
	if !err {
		return
	}
	username, err := untyped.(string)
	if !err {
		return
	}
	w.Write([]byte(username))
}

func homeHandlerGet(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func blogHandlerGet(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func contactHandlerGet(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}
