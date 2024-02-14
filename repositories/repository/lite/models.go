// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package lite

import (
	"time"
)

type Property struct {
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
	PropertyAddress     string
	PropertyTransaction string
	PropertySurface     int64
	Features            string
}

type User struct {
	ID           string
	CreatedAt    time.Time
	Isadmin      bool
	Issu         bool
	Email        string
	Passwordhash string
}
