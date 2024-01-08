package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/homepage/views/includes"
)

type HomeHandler struct {
	Log *slog.Logger
}

func NewHomeHandler(log *slog.Logger) *HomeHandler {
	return &HomeHandler{
		Log: log,
	}
}

func (hh *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got in the home function %s", r.URL.Path)
	if r.URL.Path == "/" {
		switch r.Method {
		case http.MethodGet:
			view(w, r)
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
func view(w http.ResponseWriter, r *http.Request) {
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
