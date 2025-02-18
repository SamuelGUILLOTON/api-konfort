CREATE TABLE IF NOT EXISTS users (
     id SERIAL PRIMARY KEY,
     blaze character varying NOT NULL UNIQUE,
     password_hash character varying NOT NULL,
     email character varying NOT NULL UNIQUE,
     role character varying NOT NULL,
     created_at TIMESTAMP,
     updated_at TIMESTAMP
)