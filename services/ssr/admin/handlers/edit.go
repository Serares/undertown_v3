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

type EditHandler struct {
	Log           *slog.Logger
	Service       *service.EditService
	SubmitService *service.SubmitService
}

func NewEditHandler(log *slog.Logger, service *service.EditService, submitService *service.SubmitService) *EditHandler {
	return &EditHandler{Log: log, Service: service, SubmitService: submitService}
}

func (h *EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := middleware.ID(r)
	q := r.URL.Query()

	if _, ok := q[utils.HumanReadableIdQueryKey]; ok {
		var err error
		// ❗TODO watch the delete handler
		// on reusing the token and the query params
		// this is too much spaggeti conditions
		propertyTitle := strings.Split(r.URL.Path, "/")[2]
		theId := q[utils.HumanReadableIdQueryKey][0]
		deleteUrl := fmt.Sprintf("%s/%s", types.DeletePath, utils.UrlEncodePropertyTitle(propertyTitle))
		editUrl := fmt.Sprintf("%s/%s", types.EditPath, utils.UrlEncodePropertyTitle(propertyTitle))

		// add property human readable id as query string
		deleteUrl, err = utils.AddParamToUrl(deleteUrl, utils.HumanReadableIdQueryKey, theId)
		editUrl, err = utils.AddParamToUrl(editUrl, utils.HumanReadableIdQueryKey, theId)
		if err != nil {
			h.Log.Error("malformed request", "id", theId, "error", err)
			viewEdit(w, r, types.EditProps{
				ErrorMessage:        "Failed to get the property, try again later",
				SuccessMessage:      "",
				Property:            lite.Property{},
				PropertyTypes:       types.PropertyTypes,
				PropertyTransaction: types.PropertyTransactions,
				PropertyFeatures:    types.PropertyFeatures{},
				Images:              []string{},
				FormMethod:          http.MethodPost,
				FormAction:          editUrl,
			},
				deleteUrl,
			)
			return
		}
		if r.Method == http.MethodGet {
			// populate the view with Property data
			// ❗the propertyId here is the humanReadableId
			// ❗the Features are a json string
			// should unmarshal the json string into a known struct
			// and pass the struct to the templ file to fill out the checkboxes
			property, images, propertyFeatures, err := h.Service.Get(theId, token)
			if err != nil {
				h.Log.Error("error trying to get the property", "id", theId, "error", err)
				viewEdit(w, r, types.EditProps{
					ErrorMessage:        "Failed to get the property, try again later",
					SuccessMessage:      "",
					Property:            lite.Property{},
					PropertyTypes:       types.PropertyTypes,
					PropertyTransaction: types.PropertyTransactions,
					PropertyFeatures:    types.PropertyFeatures{},
					Images:              []string{},
					FormMethod:          http.MethodPost,
					FormAction:          editUrl,
				},
					deleteUrl,
				)
				return
			}
			viewEdit(w, r, types.EditProps{
				Property:            property,
				PropertyTypes:       types.PropertyTypes,
				PropertyTransaction: types.PropertyTransactions,
				PropertyFeatures:    propertyFeatures,
				Images:              images,
				SuccessMessage:      "",
				ErrorMessage:        "",
				FormMethod:          http.MethodPost,
				FormAction:          editUrl,
			},
				deleteUrl,
			)
			return
		}

		// TODO ❗
		// This is just a patchy solution to reuse the existing code of sending data
		if r.Method == http.MethodPost {
			fullUrl := r.URL.Path + "?" + r.URL.RawQuery
			liteProperty, features, err := h.SubmitService.Submit(r, token, theId, true)
			if err != nil {
				h.Log.Error("error trying to update the property", "id", theId, "error", err)
				viewEdit(w, r, types.EditProps{
					ErrorMessage:        "Failed to send the data to server, try again later",
					SuccessMessage:      "",
					Property:            liteProperty,
					PropertyFeatures:    features,
					Images:              []string{}, // images are rip
					PropertyTypes:       types.PropertyTypes,
					PropertyTransaction: types.PropertyTransactions,
					FormMethod:          http.MethodPost,
					FormAction:          editUrl,
				},
					deleteUrl,
				)
				return
			}

			// redirect to the edit property page
			http.Redirect(w, r, fullUrl, http.StatusSeeOther)
			return
		}
	}

}

// TODO maybe I can reuse the submit.templ template
// but for now just create a new template
// The reason for the deleteUrl parm is because the types.EditProps is reused with the submit path and submit doesn't really need the delete url
func viewEdit(w http.ResponseWriter, r *http.Request, props types.EditProps, deleteUrl string) {
	views.Edit(
		types.BasicIncludes{
			Header: components.Header("Edit"),
			BannerSection: components.BannerSection(includesTypes.BannerSectionProps{
				Title: "Edit",
			},
			),
			DropzoneScript: includes.DropZone(props.Images), // TODO is this bad pattern?
			Preload:        components.Preload(),
			Navbar: components.Navbar(includesTypes.NavbarProps{
				Path:    fmt.Sprintf("admin%s", types.EditPath),
				IsAdmin: true,
			}),
			Footer:  components.Footer(),
			Scripts: components.Scripts(),
		},
		types.EditIncludes{
			HandleDeleteButton: includes.HandleDeleteButton(types.DeleteScriptProps{
				DeleteUrl: deleteUrl,
			}),
		},
		props,
	).Render(r.Context(), w)
}
