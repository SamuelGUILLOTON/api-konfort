package querys

import (
	"context"
	"errors"

	"api.konfort.com/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


type UserModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *UserModel) Insert(mail string, blaze string, password string, role models.Role) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) 
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO user (mail, blaze, password, role)
	VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(context.Background(), stmt, mail, blaze, hashedPassword, models.NEW_USER)

	if err != nil {
		return 0, err
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

	stmt := "SELECT id, password FROM user WHERE email = ?"

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

