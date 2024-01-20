package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/homepage/views/includes"
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
			viewProperty(w, r, types.BannerSectionProps{Title: ""})
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
func viewProperty(w http.ResponseWriter, r *http.Request, bannerProps types.BannerSectionProps) {
	header := includes.Header("Page title")
	preload := includes.Preload()
	navbar := includes.Navbar(types.NavbarProps{})
	footer := includes.Footer()
	scripts := includes.Scripts()
	views.Home(types.BasicIncludes{
		Header:        header,
		BannerSection: includes.BannerSection(bannerProps),
		Preload:       preload,
		Navbar:        navbar,
		Footer:        footer,
		Scripts:       scripts,
	}, types.HomeProps{
		ErrorMessage: "",
	}).Render(r.Context(), w)
}
