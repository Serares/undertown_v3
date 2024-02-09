package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/admin/middleware"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
)

type IDeleteService interface {
	Delete(humanReadableId, id, authToken string) error
}

type DeleteHandler struct {
	Log     *slog.Logger
	Service IDeleteService
}

func NewDeleteHandler(log *slog.Logger, service IDeleteService) *DeleteHandler {
	return &DeleteHandler{
		Log:     log,
		Service: service,
	}
}

// ❗ This is handled by JS/json requests
func (dh *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := middleware.ID(r)
	q := r.URL.Query()
	if _, ok := q[constants.HumanReadableIdQueryKey]; ok {
		humanReadableId := q[constants.HumanReadableIdQueryKey][0]
		// ❗the delete method is actually sent by a js script not by the form submission
		if r.Method == http.MethodDelete {
			err := dh.Service.Delete(humanReadableId, "", token)
			if err != nil {
				// respond with JSON
				dh.Log.Error("error on trying to delete the property", "err", err)
				utils.ReplyError(w, r, http.StatusInternalServerError, "Error trying to delete the property please try again later")
				return
			}
			// ❗TODO create constant variables for each path
			http.Redirect(w, r, types.ListPath, http.StatusSeeOther)
			return
		}
	}
}

// delete can be just a single_property template
// but for now I'll just use the views.Edit() component to render a delete button
func viewDelete() {

}
