CREATE TABLE IF NOT EXISTS users (
     id SERIAL PRIMARY KEY,
     blaze character varying NOT NULL UNIQUE,
     password_hash character varying NOT NULL,
     email character varying NOT NULL UNIQUE,
     role character varying NOT NULL,
     token_validation character varying,
     created_at TIMESTAMP,
     updated_at TIMESTAMP
)

CREATE TABLE email_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    subject VARCHAR(255) NOT NULL,
    html_content TEXT NOT NULL
);