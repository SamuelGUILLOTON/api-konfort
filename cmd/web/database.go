package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type role string

const (
	NEW_USER       	role = "NEW_USER"
    USER 			role = "USER"
    ADMIN     		role = "ADMIN"
	SUPER_ADMIN		role = "SUPER_ADMIN"
)

type User struct {
    ID     int64
    Mail  string
    Blaze string
    Password_hash  string
	Role  role
}

var dbpool *pgxpool.Pool

func databaseHandler() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	
	// Capture connection properties.
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

}

// users querie
func users(name string) ([]User, error) {
    // An albums slice to hold data from returned rows.
    var users []User
	
    rows, err := dbpool.Query(context.Background(), "SELECT id, mail, blaze, password_hash, role FROM users WHERE blaze = $1", name)
	if rows != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

    defer rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Mail, &user.Password_hash, &user.Role); err != nil {
            return nil, fmt.Errorf("users %v", err)
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    return albums, nil
}
