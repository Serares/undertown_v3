package handlers

import (
	"fmt"
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
		Log:           log,
		SubmitService: submitService,
	}
}

// ❗TODO
// create constants from the error messages
func (h *AdminSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			h.Log.Info("Content type", "type", r.Header.Get("Content-Type"))
			utils.ReplyError(w, r, http.StatusBadRequest, "Invalid content type")
			return
		}

		token := middleware.ID(r)
		property, features, err := h.SubmitService.Submit(r, token, "", false)

		if err != nil {
			h.Log.Error("failed to submit", err)
			viewSubmitWithFailed(w, r, types.EditProps{
				Property:            property,
				PropertyFeatures:    features,
				PropertyTypes:       types.PropertyTypes,
				PropertyTransaction: types.PropertyTransactions,
				SuccessMessage:      "",
				ErrorMessage:        "Failed to submit the property, try again",
				Images:              []string{},
				FormMethod:          http.MethodPost,
				FormAction:          types.SubmitPath,
			})
			return
		}

		viewSubmitWithFailed(w, r, types.EditProps{
			SuccessMessage:      "Success posting the property",
			ErrorMessage:        "",
			FormMethod:          http.MethodPost,
			FormAction:          types.SubmitPath,
			PropertyTypes:       types.PropertyTypes,
			PropertyTransaction: types.PropertyTransactions,
			PropertyFeatures:    types.PropertyFeatures{},
			Property:            lite.Property{},
			Images:              []string{},
		})
		return
	}
	if r.Method == http.MethodGet {
		viewSubmitWithFailed(w, r, types.EditProps{
			SuccessMessage:      "Property get success",
			ErrorMessage:        "",
			PropertyTypes:       types.PropertyTypes,
			PropertyTransaction: types.PropertyTransactions,
			FormMethod:          http.MethodPost,
			FormAction:          types.SubmitPath,
		})
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
		Preload:        components.Preload(),
		DropzoneScript: includes.DropZone(nil),
		Navbar: components.Navbar(includesTypes.NavbarProps{
			Path:    "/",
			IsAdmin: true,
		}),
		Footer:  components.Footer(),
		Scripts: components.Scripts(),
	},
		props,
	).Render(r.Context(), w)
}

// This should also support rerendering of the input values
func viewSubmitWithFailed(w http.ResponseWriter, r *http.Request, props types.EditProps) {
	views.Edit(types.BasicIncludes{
		Header: components.Header("Submit"),
		BannerSection: components.BannerSection(includesTypes.BannerSectionProps{
			Title: "Submit",
		},
		),
		Preload:        components.Preload(),
		DropzoneScript: includes.DropZone(nil),
		Navbar: components.Navbar(includesTypes.NavbarProps{
			Path:    fmt.Sprintf("admin%s", types.SubmitPath),
			IsAdmin: true,
		}),
		Footer:  components.Footer(),
		Scripts: components.Scripts(),
	},
		types.EditIncludes{}, // this is used specifically for the edit path
		props,
	).Render(r.Context(), w)
}
