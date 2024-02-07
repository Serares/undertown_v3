package types

import (
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
	"github.com/a-h/templ"
)

const (
	DeletePath   = "/delete"
	EditPath     = "/edit"
	SubmitPath   = "/submit"
	ListPath     = "/list"
	ImagesPrefix = "/assets"
	LoginPath    = "/login"
)

var DefaultMapLocation = struct {
	Lat string
	Lng string
}{
	Lat: "44.431907651074624",
	Lng: "26.097464561462406",
}

type BasicIncludes struct {
	Header        templ.Component
	BannerSection templ.Component
	Preload       templ.Component
	Navbar        templ.Component
	Footer        templ.Component
	Scripts       templ.Component
}

type EditIncludes struct {
	HandleDeleteButton templ.ComponentScript
	EditDropzoneScript templ.ComponentScript
	Modal              templ.Component
	LeafletMap         templ.ComponentScript
}

type SubmitIncludes struct {
	SubmitDropzoneScript templ.ComponentScript
	Modal                templ.Component
	LeafletMap           templ.ComponentScript
}

type SelectInputs struct {
	Value       string
	DisplayName string
}

type SubmitProps struct {
	PropertyTypes       []SelectInputs
	PropertyTransaction []SelectInputs
	PropertyFeatures    utils.PropertyFeatures
	Property            lite.Property
	SuccessMessage      string
	ErrorMessage        string
}

type EditProps struct {
	FormAction          string // needed to pass in the property human readable id
	PropertyTypes       []SelectInputs
	PropertyTransaction []SelectInputs
	PropertyFeatures    utils.PropertyFeatures
	Property            lite.Property
	ImagePaths          []string // These are the images paths, should be formatted to a path using the image names form th db
	ErrorMessage        string
	SuccessMessage      string
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

var PropertyTypes = []SelectInputs{
	{
		Value:       "APARTMENT",
		DisplayName: "Apartament",
	},
	{
		Value:       "HOUSE",
		DisplayName: "Casa",
	},
	{
		Value:       "STUDIO",
		DisplayName: "Garsoniera",
	},
	{
		Value:       "LAND",
		DisplayName: "Teren",
	},
}

var PropertyTransactions = []SelectInputs{
	{
		Value:       "0",
		DisplayName: "Vanzare",
	},
	{
		Value:       "1",
		DisplayName: "Chirie",
	},
}

type DeleteScriptProps struct {
	DeleteUrl string
}
