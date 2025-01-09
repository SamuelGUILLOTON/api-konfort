package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{s}", home)
	mux.HandleFunc("GET /snippet/view", snippetView)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("POST /snippet/create", snippetCreate)

	log.Print("starting serveur on :4000")

	err := http.ListenAndServe(":4000", mux)
	
	log.Fatal(err)
}
