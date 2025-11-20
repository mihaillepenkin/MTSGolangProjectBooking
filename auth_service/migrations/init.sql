CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (char_length(name) <= 50),
    password TEXT NOT NULL,
    email TEXT UNIQUE CHECK (char_length(email) <= 255),
    role TEXT NOT NULL CHECK (char_length(role) <= 50)
);
