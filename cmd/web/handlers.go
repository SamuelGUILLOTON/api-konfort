package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"

	"api.konfort.com/internal/models"
	"api.konfort.com/internal/querys"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Serveur", "GO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from api konfort"))
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Serveur", "GO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from api konfort"))
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	id, err := app.users.Insert(user.Email, user.Blaze, user.Password_hash, models.NEW_USER)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User %d created", id)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	
	var creds models.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	err = app.users.Authenticate(creds)

	if err != nil {
		if errors.Is(err, querys.ErrInvalidCredentials) {
			app.clientError(w, http.StatusBadRequest)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	tokenString, err := createTokenBearer(creds.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	
	w.Header().Add("Authorization", tokenString)
}

func  (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User created")
}