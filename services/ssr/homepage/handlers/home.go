package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/service"
	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/homepage/views/includes"
	"github.com/Serares/undertown_v3/utils"
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
				utils.ReplyError(w, r, http.StatusInternalServerError, "error getting the properties")
			}
			viewHome(w, r, types.HomeProps{ErrorMessage: "", FeaturedProperties: properties})
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
func viewHome(w http.ResponseWriter, r *http.Request, props types.HomeProps) {
	header := includes.Header("Page title")
	preload := includes.Preload()
	navbar := includes.Navbar()
	footer := includes.Footer()
	scripts := includes.Scripts()
	views.Home(types.HomeIncludes{
		Header:  header,
		Preload: preload,
		Navbar:  navbar,
		Footer:  footer,
		Scripts: scripts,
	}, props).Render(r.Context(), w)
}
