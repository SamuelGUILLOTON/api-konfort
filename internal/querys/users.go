package querys

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.konfort.com/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)


type UserModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *UserModel) Insert(mail string, blaze string, password string, role models.Role) (int, error) {


	now := time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) 
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (email, blaze, password_hash, role, created_at) VALUES($1, $2, $3, $4, $5)`

	result, err := m.DB.Exec(context.Background(), stmt, mail, blaze, hashedPassword, role, now)
	
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			switch pgErr.Code {
			case "23505": // Unique violation
				fmt.Println("Erreur : L'email ou le blaze existe déjà.")
			default:
				fmt.Println("Erreur SQL :", pgErr.Message)
			}
		} else {
			fmt.Println("Autre erreur :", err)
		}
	}

	insertNumber := result.RowsAffected()

	if insertNumber == 0 {
		return 0, nil
	}

	return int(insertNumber), nil
}

func (m *UserModel) Authenticate(creds models.Credentials) (error) {
	var id int
	var hashedPassword []byte


	fmt.Println("auth")	

	stmt := "SELECT id, password FROM user WHERE email = ?"

	fmt.Printf("%d %d",&id, &hashedPassword)

	err := m.DB.QueryRow(context.Background(), stmt, creds.Email).Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidCredentials
		} else {
			return err
		}
		
	} 
	
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(creds.Email))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
			} else {
			return err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return nil
}

func (m *UserModel) Exist() (bool, error) {
	return false, nil
}


// This will return a specific snippet based on its id.
func (m *UserModel) Get(id int) (models.User, error) {
	return models.User{}, nil
}

