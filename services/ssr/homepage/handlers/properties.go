package handlers

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/Serares/ssr/homepage/service"
	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/ssr/homepage/views"
	"github.com/Serares/ssr/includes/components"
)

type PropertiesHandler struct {
	Log               *slog.Logger
	PropertiesService service.PropertiesService
}

func NewPropertiesHandler(log *slog.Logger, service service.PropertiesService) *PropertiesHandler {
	return &PropertiesHandler{
		Log:               log,
		PropertiesService: service,
	}
}

func (ph *PropertiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// get the transction type from the url
		transactionType := strings.ReplaceAll(r.URL.Path, "/", "")
		bannerTitle := strings.ToUpper(transactionType)
		pagePath := r.URL.Path
		// get query strings
		sortProps := ph.getSortPropsFromQueryStrings(r.URL.Query())
		properties, err := ph.PropertiesService.ListProperties(sortProps, transactionType)
		if err != nil {
			ph.Log.Error("error getting properties", "error", err, "urlpath", r.URL.Path)
			viewProperties(w, r, types.PropertiesProps{Path: pagePath, Properties: properties, ErrorMessage: "Error fetching the properties"}, types.NavbarProps{Path: pagePath}, types.BannerSectionProps{Title: bannerTitle})
			return
		}

		viewProperties(w, r, types.PropertiesProps{Path: r.URL.Path, Properties: properties, ErrorMessage: ""}, types.NavbarProps{Path: pagePath}, types.BannerSectionProps{Title: bannerTitle})
	default:
		message := "Method not supported"
		ph.Log.Error(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
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
func viewProperties(w http.ResponseWriter, r *http.Request, props types.PropertiesProps, navbarProps types.NavbarProps, bannerProps types.BannerSectionProps) {

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
