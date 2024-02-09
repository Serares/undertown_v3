package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Serares/ssr/admin/middleware"
	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/ssr/admin/views/includes"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils"
)

type AdminSubmit struct {
	Log           *slog.Logger
	SubmitService *service.SubmitService
}

func NewSubmitHandler(log *slog.Logger, submitService *service.SubmitService) *AdminSubmit {
	return &AdminSubmit{
		Log:           log.WithGroup("Submit Handler"),
		SubmitService: submitService,
	}
}

func (h *AdminSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			h.Log.Info("Content type", "type", r.Header.Get("Content-Type"))
			utils.ReplyError(w, r, http.StatusBadRequest, "Invalid content type")
			return
		}

		token := middleware.ID(r)
		property, features, err := h.SubmitService.Submit(r, token, "")

		if err != nil {
			h.Log.Error("failed to submit", err)
			viewSubmit(w, r, types.SubmitProps{
				Property:            property,
				PropertyFeatures:    features,
				PropertyTypes:       types.PropertyTypes,
				PropertyTransaction: types.PropertyTransactions,
				SuccessMessage:      "",
				ErrorMessage:        "Failed to submit the property, try again",
			},
			)
			return
		}

		viewSubmit(w, r, types.SubmitProps{
			SuccessMessage:      "Success posting the property",
			ErrorMessage:        "",
			PropertyTypes:       types.PropertyTypes,
			PropertyTransaction: types.PropertyTransactions,
			PropertyFeatures:    utils.PropertyFeatures{},
			Property:            lite.Property{},
		},
		)
		return
	}
	if r.Method == http.MethodGet {
		viewSubmit(w, r, types.SubmitProps{
			SuccessMessage:      "",
			ErrorMessage:        "",
			PropertyTypes:       types.PropertyTypes,
			PropertyTransaction: types.PropertyTransactions,
		},
		)
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
}

// ❗ TODO is the Submit template needed?
func viewSubmit(w http.ResponseWriter, r *http.Request, props types.SubmitProps) {
	// if it failes to submit
	// reuse the edit.templ template to persist the fields in the form and
	// not have to readd the fields again if something goes wrong
	views.Submit(types.BasicIncludes{
		Header: components.Header("Submit"),
		BannerSection: components.BannerSection(includesTypes.BannerSectionProps{
			Title: "Submit",
		},
		),
		Preload: components.Preload(),
		Navbar: components.Navbar(includesTypes.NavbarProps{
			Path:    "/",
			IsAdmin: true,
		}),
		Footer:  components.Footer(),
		Scripts: components.Scripts(),
	},
		types.SubmitIncludes{
			SubmitDropzoneScript: includes.DropzoneSubmit(),
			Modal:                components.Modal(""),
			LeafletMap:           includes.LeafletMap(types.DefaultMapLocation.Lat, types.DefaultMapLocation.Lng),
		},
		props,
	).Render(r.Context(), w)
}
