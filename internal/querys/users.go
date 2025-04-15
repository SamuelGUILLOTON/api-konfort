package querys

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.konfort.com/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *UserModel) Insert(mail string, blaze string, password string, role models.Role, token_validation string) (int, error) {

	now := time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (email, blaze, password_hash, role, created_at, token_validation) VALUES($1, $2, $3, $4, $5, $6)`

	result, err := m.DB.Exec(context.Background(), stmt, mail, blaze, hashedPassword, role, now, token_validation)

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

func (m *UserModel) Authenticate(creds models.Credentials) error {
	var hashedPassword []byte
	var role string

	fmt.Println("in auth")

	stmt := "SELECT password_hash, role FROM users WHERE email = $1"

	err := m.DB.QueryRow(context.Background(), stmt, creds.Email).Scan(&hashedPassword, &role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}


	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(creds.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	if role == "NEW_USER" {
		return ErrSigninUnfinish
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

func (m *UserModel) CheckInscription(mail string) (int, error) {

	stmt := `UPDATE users SET role = 'USER' WHERE email = $1`

	result, err := m.DB.Exec(context.Background(), stmt, mail)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique violation
				return 0, fmt.Errorf("conflit de données : %v", pgErr.Message)
			default:
				return 0, fmt.Errorf("erreur SQL : %v", pgErr.Message)
			}
		}
		return 0, fmt.Errorf("erreur inconnue : %w", err)
	}

	setNumber := result.RowsAffected()

	if setNumber == 0 {
		return 0, nil
	}

	return int(setNumber), nil
}
