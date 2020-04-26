package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/personal_trivia/repo"
	"github.com/personal_trivia/routes"
)

func main() {
	log.Println("Now strarting the web server")
	repo := repo.New()
	routes := routes.New(repo)
	r := mux.NewRouter()
	// landing page
	r.HandleFunc("/", helloWorld).Methods(http.MethodGet)

	// questions
	r.HandleFunc("/question/{id}", routes.Question).Methods(http.MethodGet)
	r.HandleFunc("/question/{id}", routes.CheckAnswer).Methods(http.MethodPost)
	r.HandleFunc("/question", routes.CreateQuestion).Methods(http.MethodPost)

	// users
	r.HandleFunc("/signup", routes.CreateUser).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}
