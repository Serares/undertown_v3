package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
)

type IGetPropertyService interface {
	GetPropertyByHumanReadableId(ctx context.Context, humanReadableId string) (lite.Property, error)
	GetPropertyById(context.Context, string) (lite.Property, error)
}

type GetPropertyHandler struct {
	Log             *slog.Logger
	propertyService IGetPropertyService
}

func New(log *slog.Logger, propertyService IGetPropertyService) *GetPropertyHandler {
	return &GetPropertyHandler{
		Log:             log,
		propertyService: propertyService,
	}
}

func (gp GetPropertyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		q := r.URL.Query()
		if _, ok := q["propertyId"]; ok {
			theId := q["propertyId"][0]
			gp.Log.Info("the query param", "params", q["propertyId"][0])
			property, err := gp.propertyService.GetPropertyByHumanReadableId(r.Context(), theId)
			if err != nil {
				gp.Log.Error("error tyrying to query the property by id", "id", theId, "error", err)
				utils.ReplyError(w, r, http.StatusInternalServerError, "can't get the property")
				return
			}
			utils.ReplySuccess(w, r, http.StatusAccepted, property)
			return
		}
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "method not allowed")
}
