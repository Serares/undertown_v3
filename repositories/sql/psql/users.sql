-- name: CreateUser :exec
INSERT INTO users (
        id,
        created_at,
        isAdmin,
        isSu,
        email,
        passwordHash
    )
VALUES ($1, $2, $3, $4, $5, $6);
-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetUserByEmail :one
SELECT *
FROM users
where email = $1;
-- name: UpdateUserEmail :exec
UPDATE users
SET email = $1
WHERE id = $2;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
-- name: ListUsers :many
SELECT *
FROM users;