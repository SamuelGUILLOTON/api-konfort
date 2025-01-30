package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//le []byte() convertit la chaîne "secret-key" en slice d’octets, ce qui permet aux algorithmes de cryptographie de l'utiliser efficacement.
var secretKey = []byte("secret-key")

func createTokenBearer(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{ 
		"username": username, 
		"exp": time.Now().Add(time.Hour * 24).Unix(), 
	})

	tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 	return "Bearer " + tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return secretKey, nil
	})
   
	if err != nil {
	   return err
	}
   
	if !token.Valid {
	   return fmt.Errorf("invalid token")
	}
   
	return nil
 }