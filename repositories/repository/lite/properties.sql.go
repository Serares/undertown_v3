// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: properties.sql

package lite

import (
	"context"
	"time"
)

const addProperty = `-- name: AddProperty :exec
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
    )
`

type AddPropertyParams struct {
	ID                  string
	Humanreadableid     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Title               string
	IsProcessing        int64
	UserID              string
	Images              string
	Thumbnail           string
	IsFeatured          int64
	Price               int64
	PropertyType        string
	PropertyDescription string
	PropertyTransaction string
	PropertyAddress     string
	PropertySurface     int64
	Features            string
}

func (q *Queries) AddProperty(ctx context.Context, arg AddPropertyParams) error {
	_, err := q.db.ExecContext(ctx, addProperty,
		arg.ID,
		arg.Humanreadableid,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.IsProcessing,
		arg.UserID,
		arg.Images,
		arg.Thumbnail,
		arg.IsFeatured,
		arg.Price,
		arg.PropertyType,
		arg.PropertyDescription,
		arg.PropertyTransaction,
		arg.PropertyAddress,
		arg.PropertySurface,
		arg.Features,
	)
	return err
}

const deletePropertyByHumanReadableId = `-- name: DeletePropertyByHumanReadableId :exec
DELETE FROM properties
WHERE humanReadableId = ?
`

func (q *Queries) DeletePropertyByHumanReadableId(ctx context.Context, humanreadableid string) error {
	_, err := q.db.ExecContext(ctx, deletePropertyByHumanReadableId, humanreadableid)
	return err
}

const getByHumanReadableId = `-- name: GetByHumanReadableId :one
SELECT id, humanreadableid, created_at, updated_at, title, is_processing, user_id, images, thumbnail, is_featured, price, property_type, property_description, property_address, property_transaction, property_surface, features
FROM properties
WHERE humanReadableId = ?
LIMIT 1
`

func (q *Queries) GetByHumanReadableId(ctx context.Context, humanreadableid string) (Property, error) {
	row := q.db.QueryRowContext(ctx, getByHumanReadableId, humanreadableid)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Humanreadableid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.IsProcessing,
		&i.UserID,
		&i.Images,
		&i.Thumbnail,
		&i.IsFeatured,
		&i.Price,
		&i.PropertyType,
		&i.PropertyDescription,
		&i.PropertyAddress,
		&i.PropertyTransaction,
		&i.PropertySurface,
		&i.Features,
	)
	return i, err
}

const getProperty = `-- name: GetProperty :one
SELECT id, humanreadableid, created_at, updated_at, title, is_processing, user_id, images, thumbnail, is_featured, price, property_type, property_description, property_address, property_transaction, property_surface, features
FROM properties
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetProperty(ctx context.Context, id string) (Property, error) {
	row := q.db.QueryRowContext(ctx, getProperty, id)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Humanreadableid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.IsProcessing,
		&i.UserID,
		&i.Images,
		&i.Thumbnail,
		&i.IsFeatured,
		&i.Price,
		&i.PropertyType,
		&i.PropertyDescription,
		&i.PropertyAddress,
		&i.PropertyTransaction,
		&i.PropertySurface,
		&i.Features,
	)
	return i, err
}

const listFeaturedProperties = `-- name: ListFeaturedProperties :many
SELECT id,
    humanReadableId,
    created_at,
    title,
    thumbnail,
    price,
    property_transaction,
    property_type
FROM properties
WHERE is_featured = 1
    AND is_processing = 0
ORDER BY created_at DESC
`

type ListFeaturedPropertiesRow struct {
	ID                  string
	Humanreadableid     string
	CreatedAt           time.Time
	Title               string
	Thumbnail           string
	Price               int64
	PropertyTransaction string
	PropertyType        string
}

func (q *Queries) ListFeaturedProperties(ctx context.Context) ([]ListFeaturedPropertiesRow, error) {
	rows, err := q.db.QueryContext(ctx, listFeaturedProperties)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFeaturedPropertiesRow
	for rows.Next() {
		var i ListFeaturedPropertiesRow
		if err := rows.Scan(
			&i.ID,
			&i.Humanreadableid,
			&i.CreatedAt,
			&i.Title,
			&i.Thumbnail,
			&i.Price,
			&i.PropertyTransaction,
			&i.PropertyType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProperties = `-- name: ListProperties :many
SELECT id, humanreadableid, created_at, updated_at, title, is_processing, user_id, images, thumbnail, is_featured, price, property_type, property_description, property_address, property_transaction, property_surface, features
FROM properties
WHERE is_processing = 0
ORDER BY created_at DESC
`

func (q *Queries) ListProperties(ctx context.Context) ([]Property, error) {
	rows, err := q.db.QueryContext(ctx, listProperties)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Humanreadableid,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.IsProcessing,
			&i.UserID,
			&i.Images,
			&i.Thumbnail,
			&i.IsFeatured,
			&i.Price,
			&i.PropertyType,
			&i.PropertyDescription,
			&i.PropertyAddress,
			&i.PropertyTransaction,
			&i.PropertySurface,
			&i.Features,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPropertiesByTransactionType = `-- name: ListPropertiesByTransactionType :many
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
ORDER BY created_at DESC
`

type ListPropertiesByTransactionTypeRow struct {
	ID                  string
	Humanreadableid     string
	CreatedAt           time.Time
	Title               string
	Thumbnail           string
	Price               int64
	PropertyTransaction string
	PropertyAddress     string
	PropertySurface     int64
	Images              string
}

func (q *Queries) ListPropertiesByTransactionType(ctx context.Context, propertyTransaction string) ([]ListPropertiesByTransactionTypeRow, error) {
	rows, err := q.db.QueryContext(ctx, listPropertiesByTransactionType, propertyTransaction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPropertiesByTransactionTypeRow
	for rows.Next() {
		var i ListPropertiesByTransactionTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Humanreadableid,
			&i.CreatedAt,
			&i.Title,
			&i.Thumbnail,
			&i.Price,
			&i.PropertyTransaction,
			&i.PropertyAddress,
			&i.PropertySurface,
			&i.Images,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePropertyFeatures = `-- name: UpdatePropertyFeatures :exec
UPDATE properties
set features = ?
WHERE id = ?
`

type UpdatePropertyFeaturesParams struct {
	Features string
	ID       string
}

func (q *Queries) UpdatePropertyFeatures(ctx context.Context, arg UpdatePropertyFeaturesParams) error {
	_, err := q.db.ExecContext(ctx, updatePropertyFeatures, arg.Features, arg.ID)
	return err
}

const updatePropertyFields = `-- name: UpdatePropertyFields :exec
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
`

type UpdatePropertyFieldsParams struct {
	UpdatedAt           time.Time
	Title               string
	Images              string
	Thumbnail           string
	IsFeatured          int64
	Price               int64
	PropertyType        string
	PropertyDescription string
	PropertyTransaction string
	PropertyAddress     string
	PropertySurface     int64
	Features            string
	Humanreadableid     string
}

func (q *Queries) UpdatePropertyFields(ctx context.Context, arg UpdatePropertyFieldsParams) error {
	_, err := q.db.ExecContext(ctx, updatePropertyFields,
		arg.UpdatedAt,
		arg.Title,
		arg.Images,
		arg.Thumbnail,
		arg.IsFeatured,
		arg.Price,
		arg.PropertyType,
		arg.PropertyDescription,
		arg.PropertyTransaction,
		arg.PropertyAddress,
		arg.PropertySurface,
		arg.Features,
		arg.Humanreadableid,
	)
	return err
}

const updatePropertyImages = `-- name: UpdatePropertyImages :exec
UPDATE properties
set images = ?,
    thumbnail = ?
WHERE id = ?
`

type UpdatePropertyImagesParams struct {
	Images    string
	Thumbnail string
	ID        string
}

func (q *Queries) UpdatePropertyImages(ctx context.Context, arg UpdatePropertyImagesParams) error {
	_, err := q.db.ExecContext(ctx, updatePropertyImages, arg.Images, arg.Thumbnail, arg.ID)
	return err
}
