package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
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

	if r.Method == http.MethodGet {
		// populate the view with Property data
		q := r.URL.Query()
		// ❗the propertyId here is the humanReadableId
		if _, ok := q["propertyId"]; ok {
			theId := q["propertyId"][0]
			token, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			property, err := h.Service.Get(theId, token.Value)
			if err != nil {
				h.Log.Error("error trying to get the property", "id", theId, "error", err)
				viewEdit(w, r, types.EditProps{ErrorMessage: "Failed to get the property, try again later", SuccessMessage: ""})
				return
			}
			viewEdit(w, r, types.EditProps{
				Property:       property,
				SuccessMessage: "",
				ErrorMessage:   "",
			})
		}
	}

	if r.Method == http.MethodPost {
		// populate the view with Property data
		q := r.URL.Query()
		// ❗the propertyId here is the humanReadableId
		if _, ok := q["propertyId"]; ok {
			theId := q["propertyId"][0]
			token, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			err = h.SubmitService.Submit(r, token.Value, theId, true)
			if err != nil {
				h.Log.Error("error trying to update the property", "id", theId, "error", err)
				viewEdit(w, r, types.EditProps{ErrorMessage: "Failed to send the data to server, try again later", SuccessMessage: ""})
				return
			}

			fullUrl := r.URL.Path + "?" + r.URL.RawQuery
			// redirect to the edit proerty page
			http.Redirect(w, r, fullUrl, http.StatusSeeOther)
			return
		}
	}

}

// TODO maybe I can reuse the submit.templ template
// but for now just create a new template
func viewEdit(w http.ResponseWriter, r *http.Request, props types.EditProps) {

}
