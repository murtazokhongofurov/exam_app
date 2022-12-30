CREATE TABLE IF NOT EXISTS customers(
    id UUID PRIMARY KEY NOT NULL UNIQUE,
    full_name VARCHAR(30) NOT NULL,
    bio TEXT,
    email VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP 
);  