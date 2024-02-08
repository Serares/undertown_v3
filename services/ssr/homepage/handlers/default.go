package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	homepageViewsIncludes "github.com/Serares/ssr/homepage/views/includes"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils"
)

type ISinglePropertyService interface {
	Get(humanReadableId string) (types.ProcessedSingleProperty, error)
}

type IHomeService interface {
	Get() ([]types.ProcessedFeaturedProperty, error)
}

// ❗
// for the moment this handler will handle the home and the single property landing pages
// because the base paths are similar for now
type DefaultHandler struct {
	Log                   *slog.Logger
	HomeService           IHomeService
	SinglePropertyService ISinglePropertyService // TODO
}

func NewDefaultHandler(
	log *slog.Logger,
	homeService IHomeService,
	singlePropertyService ISinglePropertyService,
) *DefaultHandler {
	return &DefaultHandler{
		Log:                   log,
		HomeService:           homeService,
		SinglePropertyService: singlePropertyService,
	}
}

func (hh *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if _, ok := q[utils.HumanReadableIdQueryKey]; ok {
		if r.Method == http.MethodGet {
			processedProperty, err := hh.SinglePropertyService.Get(q[utils.HumanReadableIdQueryKey][0])
			path := processedProperty.Title
			if err != nil {
				ViewNotFound(w, r)
				return
			}
			viewSingleProperty(w, r,
				types.SinglePropertyViewProps{
					Property: processedProperty,
				},
				includesTypes.NavbarProps{
					Path:    path,
					IsAdmin: false,
				},
			)
			return
		}
	}

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
						Path: "/",
					})
				return
			}
			viewHome(w, r,
				types.HomeViewProps{
					ErrorMessage:       "",
					FeaturedProperties: properties,
				},
				includesTypes.NavbarProps{
					Path: "/",
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

func viewSingleProperty(w http.ResponseWriter, r *http.Request, props types.SinglePropertyViewProps, navbarProps includesTypes.NavbarProps) {

	views.Property(
		types.BasicIncludes{
			Header: components.Header("UNDERTOWN"),
			BannerSection: components.BannerSection(
				includesTypes.BannerSectionProps{
					Title: props.Property.Title,
				},
			),
			Preload: components.Preload(),
			Navbar:  components.Navbar(navbarProps),
			Footer:  components.Footer(),
			Scripts: components.Scripts(),
		},
		types.SinglePropertyIncludes{
			LeafletMap: homepageViewsIncludes.LeafletMap(
				props.Property.Features.Lat,
				props.Property.Features.Lng,
			),
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
