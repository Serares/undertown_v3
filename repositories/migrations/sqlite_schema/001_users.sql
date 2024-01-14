-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    isAdmin BOOLEAN NOT NULL,
    isSu BOOLEAN NOT NULL,
    email TEXT NOT NULL,
    -- stored as base64
    passwordHash TEXT NOT NULL
);
-- +goose Down
DROP TABLE users;