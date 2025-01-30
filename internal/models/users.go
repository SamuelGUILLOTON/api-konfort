package models

import (
	"time"
)

type Role string

const (
	NEW_USER    Role = "NEW_USER"
	USER        Role = "USER"
	ADMIN       Role = "ADMIN"
	SUPER_ADMIN Role = "SUPER_ADMIN"
)

type User struct {
	ID            int64 	`json:"id"`
	Email         string	`json:"email"`
	Blaze         string	`json:"blaze"`
	Password_hash string	`json:"password"`
	Role          Role		`json:"role"`
	created_at	  time.Time	`json:"created_at"`	
	updated_at	  time.Time `json:"updated_at"`
}

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

