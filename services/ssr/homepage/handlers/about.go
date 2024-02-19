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

type AboutHandler struct {
	Log *slog.Logger
}

func NewAboutHandler(log *slog.Logger) *AboutHandler {
	return &AboutHandler{
		Log: log,
	}
}

func (ah *AboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		viewAboutPage(
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

func viewAboutPage(w http.ResponseWriter, r *http.Request) {
	views.About(
		types.BasicIncludes{
			Header: components.Header("UNDERTOWN"),
			BannerSection: components.BannerSection(
				includesTypes.BannerSectionProps{
					Title: "ABOUT",
				},
			),
			Preload: components.Preload(),
			Navbar: components.Navbar(
				includesTypes.NavbarProps{
					Path:    constants.ABOUT_PATH,
					IsAdmin: false,
				},
			),
			Footer:  components.Footer(),
			Scripts: components.Scripts(),
		},
	).Render(r.Context(), w)
}
