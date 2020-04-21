package main

import (
	"net/http"

	"github.com/cbfield/golang-webapp-practice/models"
	"github.com/cbfield/golang-webapp-practice/routes"
	"github.com/cbfield/golang-webapp-practice/utils"
)

func main() {
	utils.LoadTemplates()
	models.Init()
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
