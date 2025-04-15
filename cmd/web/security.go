package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//le []byte() convertit la chaîne "secret-key" en slice d’octets, ce qui permet aux algorithmes de cryptographie de l'utiliser efficacement.
var secretKey = []byte("eloge-de-la-fuite")

func createBearerToken(username string) (string, error) {

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

func createTokenUrl(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{ 
		"username": username, 
		"exp": time.Now().Add(time.Hour * 3).Unix(), 
	})

	tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

	url := "https://konfort.com/user?verify=" + tokenString

	return url, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Check the signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return token, nil
}

func checkValidToken (token *jwt.Token) error {
	
	if !token.Valid {
		return fmt.Errorf("Token expiré ou invalide")
	}

	return nil
}

func checkEmailToken (token *jwt.Token) (string, error) {
	// Extraire l'email des claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("Impossible de lire le token")
	}

	email, ok := claims["email"].(string)

	if !ok {
		return "", fmt.Errorf("Impossible de lire le mail")
	}
	
	// Vérifier que c'est bien un token de vérification
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "email_verification" {
		return "", fmt.Errorf("Impossible de lire le token")
	}

	return email, nil
}