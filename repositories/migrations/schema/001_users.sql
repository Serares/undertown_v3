-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    isAdmin BOOLEAN NOT NULL,
    isSu BOOLEAN NOT NULL,
    email VARCHAR(255) NOT NULL,
    -- stored as base64
    passwordHash TEXT NOT NULL
);
-- +goose Down
DROP TABLE users;