package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
)

type PropertyHandler struct {
	Log             *slog.Logger
	PropertyHandler interface{}
}

func NewPropertyHandler(log *slog.Logger, propertyService interface{}) *PropertyHandler {
	return &PropertyHandler{
		Log: log,
	}
}

func (hh *PropertyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("property handler %s", r.URL.Path)
	if r.URL.Path == "/" {
		switch r.Method {
		case http.MethodGet:
			viewProperty(w, r, includesTypes.BannerSectionProps{Title: ""})
		default:
			message := "Method not supported"
			hh.Log.Error(message)
			http.Error(w, message, http.StatusInternalServerError)
		}
		return
	}
	// TODO handle paths that are unknown
}

// TODO should this function be defined like this?
func viewProperty(w http.ResponseWriter, r *http.Request, bannerProps includesTypes.BannerSectionProps) {
	header := components.Header("Page title")
	preload := components.Preload()
	navbar := components.Navbar(includesTypes.NavbarProps{IsAdmin: false})
	footer := components.Footer()
	scripts := components.Scripts()
	views.Home(types.BasicIncludes{
		Header:        header,
		BannerSection: components.BannerSection(bannerProps),
		Preload:       preload,
		Navbar:        navbar,
		Footer:        footer,
		Scripts:       scripts,
	}, types.HomeProps{
		ErrorMessage: "",
	}).Render(r.Context(), w)
}
