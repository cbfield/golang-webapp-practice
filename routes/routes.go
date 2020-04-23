package routes

import (
	"net/http"

	"github.com/cbfield/golang-webapp-practice/middleware"
	"github.com/cbfield/golang-webapp-practice/models"
	"github.com/cbfield/golang-webapp-practice/sessions"
	"github.com/cbfield/golang-webapp-practice/utils"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")

	r.HandleFunc("/login", loginHandlerGet).Methods("GET")
	r.HandleFunc("/login", loginHandlerPost).Methods("POST")

	r.HandleFunc("/register", registerHandlerGet).Methods("GET")
	r.HandleFunc("/register", registerHandlerPost).Methods("POST")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return r
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	updates, err := models.GetUpdates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	utils.ExecuteTemplate(w, "index.html", updates)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/", 302)
}

func loginHandlerGet(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "Username Not Recognized")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "Incorrect Password")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
	}

	userId, err := user.GetUserId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userId
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func registerHandlerGet(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerHandlerPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	err := models.RegisterUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/login", 302)
}
