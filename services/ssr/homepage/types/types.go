package types

import (
	"github.com/Serares/undertown_v3/utils"
	"github.com/a-h/templ"
)

type ProcessedSingleProperty struct {
	Title        string
	ImagePaths   []string
	DisplayPrice string
	Address      string
	Description  string
	Surface      string
	Features     utils.PropertyFeatures
}

type ProcessedFeaturedProperty struct {
	Title           string
	TransactionType string
	PropertyType    string
	PropertySurface string
	PropertyAddress string
	Price           int64
	DisplayPrice    string
	ThumbnailPath   string
	PropertyPathUrl string
	CreatedTime     string
}

// 🤔 this is for listings (chirii/vanzari)
type ProcessedListProperty struct {
	Title           string
	Address         string
	TransactionType string
	Price           int64
	DisplayPrice    string
	ThumbnailPath   string
	PropertyPathUrl string
	Surface         int64
	ImagesNumber    int64
	CreatedTime     string
}

type BasicIncludes struct {
	Header        templ.Component
	BannerSection templ.Component
	Preload       templ.Component
	Navbar        templ.Component
	Footer        templ.Component
	Scripts       templ.Component
}

type SinglePropertyIncludes struct {
	LeafletMap templ.ComponentScript
}

type HomeViewProps struct {
	ErrorMessage       string
	FeaturedProperties []ProcessedFeaturedProperty
}

type PropertiesViewProps struct {
	// TODO
	ErrorMessage string
	Path         string // it's either chirii || vanzari
	Properties   []ProcessedListProperty
}

type SinglePropertyViewProps struct {
	Property ProcessedSingleProperty // TODO define a structure for this type
}
