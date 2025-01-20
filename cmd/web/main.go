package main

import (
	"log"
	"net/http"
	"os"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dsn := os.Getenv("DATABASE_URL");

	db, err := initializeDatabase(dsn);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	

	defer db.Close()

	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /{s}", home)
	mux.HandleFunc("GET /snippet/view", snippetView)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("POST /snippet/create", snippetCreate)

	log.Print("starting serveur on :4000")

	err = http.ListenAndServe(":4000", mux)
	
	log.Fatal(err)
}

func initializeDatabase(dsn string)(*pgxpool.Pool, error){
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}