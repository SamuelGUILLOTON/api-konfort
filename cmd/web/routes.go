package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	protected := alice.New(app.requireAuthentication)

	mux.Handle("GET /{s}", protected.ThenFunc(app.home))
	mux.HandleFunc("GET /health", app.healthCheck)

	mux.HandleFunc("POST /user/create", app.userCreate)
	mux.HandleFunc("POST /user/login", app.userLogin)
	mux.HandleFunc("GET /user/verify", app.verifyUser)

	mux.HandleFunc("POST /email/create", app.emailCreate)
	mux.HandleFunc("GET /email/{name}", app.getEmailByName)

	

	//Create a middleware chain avec tout les middleware utilisé par default par notre app
	standard := alice.New(commonHeaders)

	return standard.Then(mux)
}
