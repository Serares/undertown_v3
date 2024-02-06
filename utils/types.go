package utils

import "github.com/golang-jwt/jwt/v5"

const (
	HumanReadableIdQueryKey = "ID"
	DeleteImagesFormKey     = "deletedImages"
	ImagesFormKey           = "images"
)

type JWTClaims struct {
	Email   string `json:"email"`
	UserId  string `json:"userId"`
	Isadmin bool   `json:"isadmin"`
	IsSsr   bool   `json:"isssr"`
	jwt.RegisteredClaims
}

// This is the "property" field that's received as a multipart form data POST/PUT request
type RequestProperty struct {
	Title               string                 `json:"title"`
	IsFeatured          bool                   `json:"is_featured"`
	Price               int64                  `json:"price"`
	PropertyType        string                 `json:"property_type"`
	PropertyDescription string                 `json:"property_description"`
	PropertyAddress     string                 `json:"property_address"`
	PropertyTransaction int64                  `json:"property_transaction"` // this is sent as 0 and 1 see repository/types/PropertyTransaction
	PropertySurface     int64                  `json:"property_surface"`
	Features            map[string]interface{} `json:"-"`
}
