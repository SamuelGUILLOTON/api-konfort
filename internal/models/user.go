package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type role string

const (
	NEW_USER    role = "NEW_USER"
	USER        role = "USER"
	ADMIN       role = "ADMIN"
	SUPER_ADMIN role = "SUPER_ADMIN"
)

type User struct {
	ID            int64
	Mail          string
	Blaze         string
	Password_hash string
	Role          role
}

type UserModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *UserModel) Insert(mail string, blaze string, password string, role role) (int, error) {
	stmt := `INSERT INTO user (mail, blaze, password, role)
	VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//TODO : hash password
	
	result, err := m.DB.Exec(context.Background(), stmt, mail, blaze, password, role)
	
	if err != nil {
		return 0, err
	}

	insertNumber := result.RowsAffected()

	if insertNumber == 0 {
		return 0, nil
	}

	return int(insertNumber), nil
}

// This will return a specific snippet based on its id.
func (m *UserModel) Get(id int) (User, error) {
	return User{}, nil
}

// This will return the 10 most recently created snippets.
func (m *UserModel) Latest() ([]User, error) {
	return nil, nil
}
