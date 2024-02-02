package types

import (
	"github.com/Serares/ssr/homepage/service"
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

type HomeProps struct {
	ErrorMessage       string
	FeaturedProperties []service.ProcessedFeaturedProperty
}

type ContactProps struct {
	ErrorMessage string
}

type PropertiesProps struct {
	// TODO
	ErrorMessage string
	Path         string // it's either chirii || vanzari
	Properties   []service.ProcessedProperties
}

type SinglePropertyProps struct {
	Property []lite.Property // TODO define a structure for this type
}
