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
	Log     *slog.Logger
	Service *service.EditService
}

func NewEditHandler(log *slog.Logger, service *service.EditService) *EditHandler {
	return &EditHandler{
		Log:     log.WithGroup("Edit handler"),
		Service: service,
	}
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
		deleteUrl := fmt.Sprintf("%s/%s", types.DeletePath, utils.UrlEncodeString(propertyTitle))
		editUrl := fmt.Sprintf("%s/%s", types.EditPath, utils.UrlEncodeString(propertyTitle))

		// add property human readable id as query string
		deleteUrl, err = utils.AddParamToUrl(deleteUrl, utils.HumanReadableIdQueryKey, theId)
		if err != nil {
			h.Log.Error("error creating the delete url", "error", err)
		}
		editUrl, err = utils.AddParamToUrl(editUrl, utils.HumanReadableIdQueryKey, theId)
		if err != nil {
			h.Log.Error("error creating the edit url", "error", err)
		}
		if err != nil {
			h.Log.Error("malformed request", "id", theId, "error", err)
			viewEdit(w, r, types.EditProps{
				ErrorMessage:        "Failed to get the property, try again later",
				SuccessMessage:      "",
				Property:            lite.Property{},
				PropertyTypes:       types.PropertyTypes,
				PropertyTransaction: types.PropertyTransactions,
				PropertyFeatures:    utils.PropertyFeatures{},
				ImagePaths:          []string{},
				FormAction:          editUrl,
			},
				deleteUrl,
				http.StatusInternalServerError,
			)
			return
		}
		if r.Method == http.MethodGet {
			// populate the view with Property data
			// ❗the propertyId here is the humanReadableId
			// ❗the Features are a json string
			// should unmarshal the json string into a known struct
			// and pass the struct to the templ file to fill out the checkboxes
			property, imagesPaths, propertyFeatures, err := h.Service.Get(theId, token)
			if err != nil {
				h.Log.Error("error trying to get the property", "id", theId, "error", err)
				viewEdit(w, r, types.EditProps{
					ErrorMessage:        "Failed to get the property, try again later",
					SuccessMessage:      "",
					Property:            lite.Property{},
					PropertyTypes:       types.PropertyTypes,
					PropertyTransaction: types.PropertyTransactions,
					PropertyFeatures:    utils.PropertyFeatures{},
					ImagePaths:          []string{},
					FormAction:          editUrl,
				},
					deleteUrl,
					http.StatusInternalServerError,
				)
				return
			}
			viewEdit(w, r, types.EditProps{
				Property:            property,
				PropertyTypes:       types.PropertyTypes,
				PropertyTransaction: types.PropertyTransactions,
				PropertyFeatures:    propertyFeatures,
				ImagePaths:          imagesPaths,
				SuccessMessage:      "",
				ErrorMessage:        "",
				FormAction:          editUrl,
			},
				deleteUrl,
				http.StatusOK,
			)
			return
		}

		// TODO ❗
		// This is just a patchy solution to reuse the existing code of sending data
		if r.Method == http.MethodPost {
			liteProperty, features, err := h.Service.Post(r, token, theId)
			if err != nil {
				h.Log.Error("error trying to update the property", "id", theId, "error", err)
				viewEdit(w, r, types.EditProps{
					ErrorMessage:        "Failed to send the data to server, try again later",
					SuccessMessage:      "",
					Property:            liteProperty,
					PropertyFeatures:    features,
					ImagePaths:          []string{}, // images are rip
					PropertyTypes:       types.PropertyTypes,
					PropertyTransaction: types.PropertyTransactions,
					FormAction:          editUrl,
				},
					deleteUrl,
					http.StatusInternalServerError,
				)
				return
			}

			// ❗TODO
			// If the user changes the property title, the ridirect will display the old property title in the url path
			// because on a success backend query the liteProperty, features, err are all nullish values
			fullUrl := r.URL.Path
			fullUrl, err = utils.AddParamToUrl(fullUrl, utils.HumanReadableIdQueryKey, q[utils.HumanReadableIdQueryKey][0])
			if err != nil {
				h.Log.Error("error trying to create the redirect url", "error", err)
				http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, fullUrl, http.StatusSeeOther)
			return
		}
	}

}

// TODO maybe I can reuse the submit.templ template
// but for now just create a new template
// The reason for the deleteUrl parm is because the types.EditProps is reused with the submit path and submit doesn't really need the delete url
func viewEdit(w http.ResponseWriter, r *http.Request, props types.EditProps, deleteUrl string, statusCode int64) {
	w.WriteHeader(int(statusCode))
	views.Edit(
		types.BasicIncludes{
			Header: components.Header("Edit"),
			BannerSection: components.BannerSection(includesTypes.BannerSectionProps{
				Title: "Edit",
			},
			),
			Preload: components.Preload(),
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
			EditDropzoneScript: includes.DropzoneEdit(props.ImagePaths, props.FormAction, utils.DeleteImagesFormKey),
			Modal:              components.Modal(""),
			LeafletMap:         includes.LeafletMap(props.PropertyFeatures.Lat, props.PropertyFeatures.Lng),
		},
		props,
	).Render(r.Context(), w)
}
