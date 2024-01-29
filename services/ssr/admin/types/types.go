package types

import "github.com/a-h/templ"

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
