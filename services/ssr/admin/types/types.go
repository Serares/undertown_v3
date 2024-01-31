package types

import (
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/a-h/templ"
)

type BasicIncludes struct {
	Header        templ.Component
	BannerSection templ.Component
	Preload       templ.Component
	Navbar        templ.Component
	Footer        templ.Component
	Scripts       templ.Component
}

type SubmitProps struct {
	Message      string
	ErrorMessage string
}

type EditProps struct {
	Property       lite.Property
	ErrorMessage   string
	SuccessMessage string
}

type LoginProps struct {
	Message      string
	ErrorMessage string
}

type ListingProperty struct {
	Title              string
	Address            string
	TransactionType    string
	Price              int64
	DisplayPrice       string
	Thumbnail          string
	EditPropertyPath   string
	DeletePropertyPath string
	Surface            int64
	ImagesNumber       int64
	CreatedTime        string
}

type ListingProps struct {
	Properties     []ListingProperty
	ErrorMessage   string
	SuccessMessage string
}
