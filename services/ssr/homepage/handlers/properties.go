package handlers

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/Serares/ssr/homepage/service"
	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	homepageViewsIncludes "github.com/Serares/ssr/homepage/views/includes"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils/constants"
)

type ISinglePropertyService interface {
	Get(humanReadableId string) (types.ProcessedSingleProperty, error)
}

type PropertiesHandler struct {
	Log                   *slog.Logger
	PropertiesService     service.PropertiesService
	SinglePropertyService ISinglePropertyService
}

func NewPropertiesHandler(
	log *slog.Logger,
	service service.PropertiesService,
	singlePropertyService ISinglePropertyService,
) *PropertiesHandler {
	return &PropertiesHandler{
		Log:                   log.WithGroup("Properties Handler"),
		PropertiesService:     service,
		SinglePropertyService: singlePropertyService,
	}
}

func (ph *PropertiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	if _, ok := q[constants.HumanReadableIdQueryKey]; ok {
		if r.Method == http.MethodGet {
			processedProperty, err := ph.SinglePropertyService.Get(q[constants.HumanReadableIdQueryKey][0])
			path := r.URL.Path
			if err != nil {
				ph.Log.Error("error trying to render the Single property", "error", err)
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
		// get the transction type from the url
		translatedTransactionType := strings.ReplaceAll(r.URL.Path, "/", "")
		bannerTitle := strings.ToUpper(translatedTransactionType)
		pagePath := r.URL.Path
		// get query strings
		sortProps := ph.getSortPropsFromQueryStrings(r.URL.Query())
		properties, err := ph.PropertiesService.ListProperties(sortProps, translatedTransactionType)
		if err != nil {
			ph.Log.Error("error getting properties", "error", err, "urlpath", r.URL.Path)
			viewProperties(
				w,
				r,
				types.PropertiesViewProps{
					Path:         pagePath,
					Properties:   properties,
					ErrorMessage: "Error fetching the properties",
				},
				includesTypes.NavbarProps{
					Path: pagePath,
				},
				includesTypes.BannerSectionProps{
					Title: bannerTitle,
				})
			return
		}

		viewProperties(
			w,
			r,
			types.PropertiesViewProps{
				Path:         r.URL.Path,
				Properties:   properties,
				ErrorMessage: "",
			},
			includesTypes.NavbarProps{
				Path: pagePath,
			},
			includesTypes.BannerSectionProps{
				Title: bannerTitle,
			})
		return
	}

	message := "Method not supported"
	ph.Log.Error(message)
	http.Error(w, message, http.StatusMethodNotAllowed)
	return
}

// TODO test this function
func (ph *PropertiesHandler) getSortPropsFromQueryStrings(queryStrings url.Values) service.SortProps {
	ph.Log.Info("those are the queries", "queries", queryStrings)
	var sortProps service.SortProps
	sortOrder := queryStrings.Get("sort_order")
	if sortOrder != "" {
		values := strings.Split(sortOrder, "/")
		switch values[0] {
		case "price":
			sortProps.Price = resolveClientSortDirection(values[1])
		case "surface":
			sortProps.Surface = resolveClientSortDirection(values[1])
		case "createdAt":
			sortProps.PublishedDate = resolveClientSortDirection(values[1])
		}
	}
	return service.SortProps{}
}

func resolveClientSortDirection(sortDirection string) string {
	switch sortDirection {
	case "desc":
		return service.DSC
	case "asc":
		return service.ASC
	default:
		return ""
	}
}

// TODO should this function be defined like this?
func viewProperties(
	w http.ResponseWriter,
	r *http.Request,
	props types.PropertiesViewProps,
	navbarProps includesTypes.NavbarProps,
	bannerProps includesTypes.BannerSectionProps,
) {

	header := components.Header(bannerProps.Title)
	preload := components.Preload()
	navbar := components.Navbar(navbarProps)
	footer := components.Footer()
	scripts := components.Scripts()
	views.Properties(types.BasicIncludes{
		Header:        header,
		Preload:       preload,
		BannerSection: components.BannerSection(bannerProps),
		Navbar:        navbar,
		Footer:        footer,
		Scripts:       scripts,
	}, props).Render(r.Context(), w)
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
