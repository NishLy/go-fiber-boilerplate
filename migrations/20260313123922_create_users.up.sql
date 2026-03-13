CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE,
    password TEXT,
    created_at TIMESTAMP DEFAULT now()
);