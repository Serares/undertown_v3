package types

import (
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/a-h/templ"
)

type HomeIncludes struct {
	Header  templ.Component
	Preload templ.Component
	Navbar  templ.Component
	Footer  templ.Component
	Scripts templ.Component
}

type HomeProps struct {
	ErrorMessage       string
	FeaturedProperties []lite.ListFeaturedPropertiesRow
}

type ContactIncludes struct {
	Header  templ.Component
	Preload templ.Component
	Navbar  templ.Component
	Footer  templ.Component
	Scripts templ.Component
}

type ContactProps struct {
	ErrorMessage string
}
