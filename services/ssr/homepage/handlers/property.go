package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/service"
	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/homepage/views/includes"
)

type PropertyHandler struct {
	Log         *slog.Logger
	HomeService service.HomeService
}

func NewPropertyHandler(log *slog.Logger, homeService service.HomeService) *HomeHandler {
	return &HomeHandler{
		Log:         log,
		HomeService: homeService,
	}
}

func (hh *PropertyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("property handler %s", r.URL.Path)
	if r.URL.Path == "/" {
		switch r.Method {
		case http.MethodGet:
			viewProperty(w, r)
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
func viewProperty(w http.ResponseWriter, r *http.Request) {
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
	}, types.HomeProps{
		ErrorMessage: "",
	}).Render(r.Context(), w)
}
