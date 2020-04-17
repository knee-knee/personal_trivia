package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Now strarting the web server")
	http.HandleFunc("/", helloWorld)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}
