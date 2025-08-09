package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"api.konfort.com/internal/models"
	"api.konfort.com/internal/querys"
	"api.konfort.com/internal/services"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Serveur", "GO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from api konfort"))
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Serveur", "GO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from KONFORT"))
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	
	//fetch body request
	decoder := json.NewDecoder(r.Body)
	var user models.UserPost
	
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	fmt.Println(user)
	
	//create token 
	tokenString, err := createTokenUrl(user.Email)

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	
	//save user in bdd
	id, err := app.users.Insert(user.Email, user.Blaze, user.Password_hash, models.NEW_USER, tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// if one user has been added it has been created
	if id == 1 {
		fmt.Println("tettetetetetet")
		var emailTemplate *models.Email

		var template string = "check_sign_up";

		//fetch emailTemplate
		emailTemplate, err = app.email.GetEmailTemplate(template)

		if err != nil {
			app.logger.Error("Failed to get email template", "error", err)
			return
		}

		//send mail to user
		err := services.Mailer("samuel.guilloton01@gmail.com", "samuel.guilloton01@gmail.com", 
		emailTemplate.Subject, tokenString, emailTemplate.Html_content)
	
		if err != nil {
			fmt.Printf("Email sending error: %v\n", err)
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}
		
		// user has been created in bdd partialy and email has been send
		w.WriteHeader(http.StatusCreated)
		return
	}
	
	// If the user is created successfully but email is not sent
	w.WriteHeader(http.StatusOK)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inHandler")

	var creds models.Credentials
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&creds)

	fmt.Printf("%s, %s", creds.Email, creds.Password)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	err = app.users.Authenticate(creds)
	fmt.Printf("%d", err)

	if err != nil {
		if errors.Is(err, querys.ErrInvalidCredentials) {
			app.clientError(w, http.StatusBadRequest)
		}
		if errors.Is(err, querys.ErrSigninUnfinish) {
			app.clientError(w, http.StatusBadRequest)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	tokenString, err := createBearerToken(creds.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Authorization", tokenString)
}

func (app *application) emailCreate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var email models.Email

	err := decoder.Decode(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := app.email.InsertEmailTemplate(email.Name, email.Subject, email.Html_content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%d created", id)
	w.WriteHeader(http.StatusOK)
}

func (app *application) getEmailByName(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")

	if name == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email, err := app.email.GetEmailTemplate(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(email)
	w.WriteHeader(http.StatusOK)
}

// Route de vérification
func (app *application) verifyUser(w http.ResponseWriter, r *http.Request) {
	// Récupérer le token depuis l'URL
	tokenString := r.PathValue("token")

	token, err := verifyToken(tokenString);

	if err != nil || token == nil{
		http.Error(w, "Token invalide", http.StatusBadRequest)
		return
	}

	err = checkValidToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := checkEmailToken(token)

	
	if err != nil || email == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mettre à jour le statut de l'utilisateur
	rowsAffected, err := app.users.CheckInscription(email)

	if err != nil {
		http.Error(w, "Erreur de vérification", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Aucun utilisateur mis à jour. Vérifiez l'adresse email.", http.StatusNotFound)
		return
	}

	// Répondre avec succès
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Votre email a été vérifié avec succès"))
}


