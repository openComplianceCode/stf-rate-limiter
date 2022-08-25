package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func main() {
	r := chi.NewRouter()
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8001", nil))
}
