package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	protected := alice.New(app.requireAuthentication)

	mux.HandleFunc("GET /{s}", app.home)
	mux.HandleFunc("GET /health", app.healthCheck)
	mux.Handle("POST /user/create", protected.ThenFunc(app.userCreate))
	mux.HandleFunc("POST /user/login", app.userLogin)
	mux.HandleFunc("POST /user/logout", app.userLogout)

	//Create a middleware chain avec tout les middleware utilisé par default par notre app
	standard := alice.New(app.recoverPanic, commonHeaders)

	return standard.Then(mux)
}
