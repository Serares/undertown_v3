-- name: CreateUser :exec
INSERT INTO users (
        id,
        created_at,
        isAdmin,
        isSu,
        email,
        passwordHash
    )
VALUES (?, ?, ?, ?, ?, ?);
-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?;
-- name: GetUserByEmail :one
SELECT *
FROM users
where email = ?;
-- name: UpdateUserEmail :exec
UPDATE users
SET email = ?
WHERE id = ?;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
-- name: ListUsers :many
SELECT *
FROM users;