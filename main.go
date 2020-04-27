package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/personal_trivia/middleware"
	"github.com/personal_trivia/repo"
	"github.com/personal_trivia/routes"
)

func main() {
	log.Println("Now strarting the web server")
	repo := repo.New()
	routes := routes.New(repo)
	mw := middleware.New(repo)
	r := mux.NewRouter()
	// landing page
	r.HandleFunc("/", healthCheck).Methods(http.MethodGet)

	// questions
	r.Handle("/question/{id}", mw.AuthMiddleware(routes.Question)).Methods(http.MethodGet)
	r.Handle("/question/{id}", mw.AuthMiddleware(routes.CheckAnswer)).Methods(http.MethodPost)
	r.Handle("/question", mw.AuthMiddleware(routes.CreateQuestion)).Methods(http.MethodPost)

	// users
	r.HandleFunc("/signup", routes.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", routes.Login).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "healthy")
}
