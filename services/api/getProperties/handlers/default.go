package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
)

type FeaturedPropertiesResponse struct {
	Results []lite.ListFeaturedPropertiesRow
}

// TODO the property has to be striped of some columns (user_id and id)
type PropertiesResponse struct {
	Results []lite.Property
}

type IGetPropertiesService interface {
	ListFeaturedProperties(context.Context) (*[]lite.ListFeaturedPropertiesRow, error)
	ListProperties(context.Context) (*[]lite.Property, error)
}

type GetPropertiesHandler struct {
	Log                  *slog.Logger
	GetPropertiesService IGetPropertiesService
}

func New(log *slog.Logger, getPropertiesService IGetPropertiesService) *GetPropertiesHandler {
	return &GetPropertiesHandler{
		Log:                  log,
		GetPropertiesService: getPropertiesService,
	}
}

// ❗
// TODO
// try to refactor and not repeat the code this much
func (gh GetPropertiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := r.URL.Query()

		if _, ok := q["featured"]; ok {
			properties, err := gh.GetPropertiesService.ListFeaturedProperties(r.Context())
			if err != nil {
				gh.Log.Error("error trying to get featured properties", "error", err)
				utils.ReplyError(w, r, http.StatusInternalServerError, fmt.Sprintf("error trying to get the featured properties %v", err))
				return
			}
			response := FeaturedPropertiesResponse{
				Results: *properties,
			}

			utils.ReplySuccess(w, r, http.StatusOK, response)
			return
		}

		properties, err := gh.GetPropertiesService.ListProperties(r.Context())
		if err != nil {
			gh.Log.Error("error trying to get featured properties", "error", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, fmt.Sprintf("error trying to get the featured properties %v", err))
			return
		}
		response := PropertiesResponse{
			Results: *properties,
		}

		utils.ReplySuccess(w, r, http.StatusOK, response)
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "")
	return
}
