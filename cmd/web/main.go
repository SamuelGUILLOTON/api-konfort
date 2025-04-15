package main

import (
	"log"
	"net/http"
	"log/slog"
	"os"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"api.konfort.com/internal/querys"
)


// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger *slog.Logger
	users  *querys.UserModel
	email  *querys.MailModel
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := os.Getenv("DATABASE_URL");

	db, err := initializeDatabase(dsn);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	
	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
		users: &querys.UserModel{DB: db},
		email: &querys.MailModel{DB: db},
	}
	defer db.Close()

	log.Print("starting serveur on :4000")
	err = http.ListenAndServe(":4000", app.routes())
	log.Fatal(err)
}

func initializeDatabase(dsn string)(*pgxpool.Pool, error){
	dbpool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}