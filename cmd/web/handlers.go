package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Serveur", "GO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from api konfort"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	snippetId, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || snippetId < 1 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Display a specific snippet with id %d", snippetId)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Serveur", "GO")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Display a form for creating a new snippet..."))
}
