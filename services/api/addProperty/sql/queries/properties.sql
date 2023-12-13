-- name: AddProperty :exec
INSERT INTO properties(
        id,
        created_at,
        updated_at,
        title,
        floor,
        published_at,
        user_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7);