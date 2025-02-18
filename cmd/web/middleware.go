package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
	"encoding/json"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {

		defer func()  {
			if err := recover(); err != nil {
				// quand go voit se header dans la rep il ferme ensuite la connexion
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// If the user is not authenticated, redirect them to the login page and
	// return from the middleware chain so that no subsequent handlers in
	// the chain are executed.
	Bearer := r.Header.Get("Authorization")

	err := verifyToken(Bearer)
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	next.ServeHTTP(w, r)
	})
}


func (app *application) loggingPostRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		var err error

		if r.Body != nil {
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("Body reading error: %v", err)
				return
			}
			defer r.Body.Close()
		}

		fmt.Printf("Headers: %+v\n", r.Header)

		var prettyJSON bytes.Buffer
		
		if len(bodyBytes) > 0 {
			
			if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
				fmt.Printf("JSON parse error: %v", err)
				return
			}
			fmt.Println(string(prettyJSON.String()))
		} else {
			fmt.Printf("Body: No Body Supplied\n")
		}

		next.ServeHTTP(w, r)
	})
}