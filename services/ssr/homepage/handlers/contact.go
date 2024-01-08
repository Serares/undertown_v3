package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/homepage/views/includes"
)

type ContactHandler struct {
	Log *slog.Logger
}

func NewContactHandler(log *slog.Logger) *ContactHandler {
	return &ContactHandler{
		Log: log,
	}
}

func (ch *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got in the contact function %s", r.URL.Path)
	if r.URL.Path == "" {
		switch r.Method {
		case http.MethodGet:
			viewContact(w, r)
		default:
			message := "Method not supported"
			ch.Log.Error(message)
			http.Error(w, message, http.StatusInternalServerError)
		}
		return
	}
	// TODO handle paths that are unknown
}

// TODO should this function be defined like this?
func viewContact(w http.ResponseWriter, r *http.Request) {
	header := includes.Header("Page title")
	preload := includes.Preload()
	navbar := includes.Navbar()
	footer := includes.Footer()
	scripts := includes.Scripts()
	views.Contact(types.ContactIncludes{
		Header:  header,
		Preload: preload,
		Navbar:  navbar,
		Footer:  footer,
		Scripts: scripts,
	}, types.ContactProps{
		ErrorMessage: "",
	}).Render(r.Context(), w)
}
