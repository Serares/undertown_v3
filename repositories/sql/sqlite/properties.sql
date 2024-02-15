-- name: AddProperty :exec
INSERT INTO properties(
        id,
        humanReadableId,
        created_at,
        updated_at,
        title,
        is_processing,
        user_id,
        images,
        thumbnail,
        is_featured,
        price,
        property_type,
        property_description,
        property_transaction,
        property_address,
        property_surface,
        features
    )
VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    );
-- name: GetProperty :one
SELECT *
FROM properties
WHERE id = ?
LIMIT 1;
-- name: GetByHumanReadableId :one
SELECT *
FROM properties
WHERE humanReadableId = ?
LIMIT 1;
-- name: DeletePropertyByHumanReadableId :exec
DELETE FROM properties
WHERE humanReadableId = ?;
-- name: ListProperties :many
SELECT *
FROM properties
WHERE is_processing = 0
ORDER BY created_at DESC;
-- name: ListFeaturedProperties :many
SELECT id,
    humanReadableId,
    created_at,
    title,
    thumbnail,
    price,
    property_type,
    property_transaction,
    property_address,
    property_surface,
    images
FROM properties
WHERE is_featured = 1
    AND is_processing = 0
ORDER BY created_at DESC;
-- name: ListPropertiesByTransactionType :many
SELECT id,
    humanReadableId,
    created_at,
    title,
    thumbnail,
    price,
    property_transaction,
    property_address,
    property_surface,
    images
FROM properties
WHERE property_transaction = ?
    AND is_processing = 0
ORDER BY created_at DESC;
-- name: UpdatePropertyImages :exec
UPDATE properties
set images = ?,
    thumbnail = ?
WHERE id = ?;
-- name: UpdatePropertyFeatures :exec
UPDATE properties
set features = ?
WHERE id = ?;
-- name: UpdatePropertyFields :exec
UPDATE properties
set updated_at = ?,
    title = ?,
    images = ?,
    thumbnail = ?,
    is_featured = ?,
    price = ?,
    property_type = ?,
    property_description = ?,
    property_transaction = ?,
    property_address = ?,
    property_surface = ?,
    features = ?
WHERE humanReadableId = ?