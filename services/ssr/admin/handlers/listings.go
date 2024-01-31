package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
)

type AdminListings struct {
	Log     *slog.Logger
	Service *service.ListingsService
}

func NewListingsHandler(log *slog.Logger, service *service.ListingsService) *AdminListings {
	return &AdminListings{
		Log:     log,
		Service: service,
	}
}

func (h *AdminListings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("token")
		if err != nil {
			viewListings(w, r, types.ListingProps{Properties: []types.ListingProperty{}, ErrorMessage: "Invalid authentication", SuccessMessage: ""})
			return
		}
		properties, err := h.Service.List(cookie.Value)
		if err != nil {
			h.Log.Error("error invalid response", err, err)
			viewListings(w, r, types.ListingProps{
				Properties:     []types.ListingProperty{},
				ErrorMessage:   "invalid email or password",
				SuccessMessage: "",
			})
			return
		}

		viewListings(w, r, types.ListingProps{
			Properties:     properties,
			SuccessMessage: "Success getting properties",
			ErrorMessage:   "",
		})
		return
	}
}

func viewListings(w http.ResponseWriter, r *http.Request, props types.ListingProps) {
	views.Listings(
		types.BasicIncludes{
			Header:        components.Header("Login"),
			BannerSection: components.BannerSection(includesTypes.BannerSectionProps{Title: "Login"}),
			Preload:       components.Preload(),
			Navbar:        components.Navbar(includesTypes.NavbarProps{Path: "/listings", IsAdmin: true}),
			Footer:        components.Footer(),
			Scripts:       components.Scripts(),
		},
		props,
	).Render(r.Context(), w)
}
