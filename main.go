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
	fmt.Println("Now strarting the web server")
	repo := repo.New()
	routes := routes.New(repo)
	r := mux.NewRouter()
	r.HandleFunc("/", helloWorld).Methods(http.MethodGet)
	r.HandleFunc("/questions/{id}", routes.Question).Methods(http.MethodGet)
	r.HandleFunc("/question", routes.CreateQuestion)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}
