package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
)

type IHomeService interface {
	Get() ([]types.ProcessedFeaturedProperty, error)
}

// ❗
// for the moment this handler will handle the home and the single property landing pages
// because the base paths are similar for now
type DefaultHandler struct {
	Log         *slog.Logger
	HomeService IHomeService
}

func NewDefaultHandler(
	log *slog.Logger,
	homeService IHomeService,
) *DefaultHandler {
	return &DefaultHandler{
		Log:         log,
		HomeService: homeService,
	}
}

func (hh *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		if r.URL.Path == "/" {
			properties, err := hh.HomeService.Get()
			if err != nil {
				hh.Log.Error("error getting properties", "error", err)
				viewHome(w, r,
					types.HomeViewProps{
						ErrorMessage:       "Error getting the properties",
						FeaturedProperties: properties,
					},
					includesTypes.NavbarProps{
						Path: "/home",
					})
				return
			}
			viewHome(w, r,
				types.HomeViewProps{
					ErrorMessage:       "",
					FeaturedProperties: properties,
				},
				includesTypes.NavbarProps{
					Path: "/home",
				})
			return
		}
	}
	ViewNotFound(w, r)
}

// TODO handle paths that are unknown

// TODO should this function be defined like this?
func viewHome(
	w http.ResponseWriter,
	r *http.Request,
	props types.HomeViewProps,
	navbarProps includesTypes.NavbarProps,
) {
	views.Home(
		types.BasicIncludes{
			Header:  components.Header("UNDERTOWN"),
			Preload: components.Preload(),
			Navbar:  components.Navbar(navbarProps),
			Footer:  components.Footer(),
			Scripts: components.Scripts(),
		},
		props,
	).Render(r.Context(), w)
}

func ViewNotFound(w http.ResponseWriter, r *http.Request) {
	views.NotFound(
		types.BasicIncludes{
			Header: components.Header("UNDERTOWN"),
			BannerSection: components.BannerSection(
				includesTypes.BannerSectionProps{
					Title: "Not Found",
				},
			),
			Preload: components.Preload(),
			Navbar: components.Navbar(
				includesTypes.NavbarProps{
					Path:    "OOPS",
					IsAdmin: false,
				},
			),
			Footer:  components.Footer(),
			Scripts: components.Scripts(),
		},
	).Render(r.Context(), w)
}
