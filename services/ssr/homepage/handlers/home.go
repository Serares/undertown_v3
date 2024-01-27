package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/service"
	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/includes/components"
)

type HomeHandler struct {
	Log         *slog.Logger
	HomeService service.HomeService
}

func NewHomeHandler(log *slog.Logger, homeService service.HomeService) *HomeHandler {
	return &HomeHandler{
		Log:         log,
		HomeService: homeService,
	}
}

func (hh *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got in the home function %s", r.URL.Path)
	if r.URL.Path == "/" {
		switch r.Method {
		case http.MethodGet:
			properties, err := hh.HomeService.ListProperties()
			if err != nil {
				hh.Log.Error("error getting properties", "error", err)
				viewHome(w, r, types.HomeProps{ErrorMessage: "Error getting the properties", FeaturedProperties: properties}, types.NavbarProps{Path: "/"})
				return
			}
			viewHome(w, r, types.HomeProps{ErrorMessage: "", FeaturedProperties: properties}, types.NavbarProps{Path: "/"})
			return
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
func viewHome(w http.ResponseWriter, r *http.Request, props types.HomeProps, navbarProps types.NavbarProps) {
	header := components.Header("UNDERTOWN")
	preload := components.Preload()
	navbar := components.Navbar(navbarProps)
	footer := components.Footer()
	scripts := components.Scripts()
	views.Home(types.BasicIncludes{
		Header:  header,
		Preload: preload,
		Navbar:  navbar,
		Footer:  footer,
		Scripts: scripts,
	}, props).Render(r.Context(), w)
}
