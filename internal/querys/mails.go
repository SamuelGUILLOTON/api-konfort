package querys

import (
	"context"
	"errors"
	"fmt"

	"api.konfort.com/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)


type MailModel struct {
	DB *pgxpool.Pool
}

func (m *MailModel) InsertEmailTemplate(name string, subject string, html_content string) (int, error) {
	var id int

	stmt := `INSERT INTO email_templates (name, subject, html_content) VALUES ($1, $2, $3) RETURNING id`

	err := m.DB.QueryRow(context.Background(), stmt, name, subject, html_content).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique constraint violation
				return 0, fmt.Errorf("erreur : le nom existe déjà")
			default:
				return 0, fmt.Errorf("erreur SQL : %s", pgErr.Message)
			}
		}
		return 0, fmt.Errorf("autre erreur : %w", err)
	}

	return id, nil
}


func (m *MailModel) GetEmailTemplate(name string) (*models.Email, error) {
	
	fmt.Printf("%s", name);
	var email models.Email

	stmt := "SELECT name, subject, html_content FROM email_templates WHERE name = $1"

	err := m.DB.QueryRow(context.Background(), stmt, name).Scan(&email.Name, &email.Subject, &email.Html_content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("aucun template trouvé pour '%s'", name)
		}
		return nil, fmt.Errorf("erreur SQL : %w", err)
	}

	return &email, nil
}
