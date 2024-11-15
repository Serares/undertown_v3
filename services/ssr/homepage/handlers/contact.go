package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
)

type ContactHandler struct {
	Log *slog.Logger
}

func NewContactHandler(log *slog.Logger) *ContactHandler {
	return &ContactHandler{
		Log: log,
	}
}

func (ah *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		viewContactPage(
			w,
			r,
		)
		return
	}
	utils.ReplyError(
		w,
		r,
		http.StatusMethodNotAllowed,
		"Method not allowed",
	)
}

func viewContactPage(w http.ResponseWriter, r *http.Request) {
	views.Contact(
		types.BasicIncludes{
			Header: components.Header("UNDERTOWN"),
			BannerSection: components.BannerSection(
				includesTypes.BannerSectionProps{
					Title: "CONTACT",
				},
			),
			Preload: components.Preload(),
			Navbar: components.Navbar(
				includesTypes.NavbarProps{
					Path:    constants.CONTACT_PATH,
					IsAdmin: false,
				},
			),
			Footer:  components.Footer(),
			Scripts: components.Scripts(),
		},
	).Render(r.Context(), w)
}
