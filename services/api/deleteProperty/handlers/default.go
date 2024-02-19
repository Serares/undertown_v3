package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
)

type DeleteResponse struct {
	Results []string
}

type IDeletePropertyService interface {
	DeleteProperty(ctx context.Context, id, humanReadableId string) error
}

type DeletePropertyHandler struct {
	Log     *slog.Logger
	Service IDeletePropertyService
}

func New(log *slog.Logger, service IDeletePropertyService) *DeletePropertyHandler {
	return &DeletePropertyHandler{
		Log:     log,
		Service: service,
	}
}

// ❗
// TODO
// try to refactor and not repeat the code this much
func (gh DeletePropertyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		q := r.URL.Query()
		// ❗TODO at this point it will make sense to create a midleware util module to get the query params/cookies
		if _, ok := q[constants.QUERY_PARAMETER_HUMANREADABLEID]; ok {
			theId := q[constants.QUERY_PARAMETER_HUMANREADABLEID][0]
			err := gh.Service.DeleteProperty(r.Context(), "", theId)
			if err != nil {
				gh.Log.Error("error trying to delete the property", "error", err)
				utils.ReplyError(w, r, http.StatusInternalServerError, fmt.Sprintf("error trying to delete the property %v : %s", err, theId))
				return
			}
			response := DeleteResponse{
				Results: []string{"Success delete"},
			}

			utils.ReplySuccess(w, r, http.StatusOK, response)
			return
		}
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "")
}
