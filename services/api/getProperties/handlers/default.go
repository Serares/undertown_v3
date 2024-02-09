package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
)

type FeaturedPropertiesResponse struct {
	Results []lite.ListFeaturedPropertiesRow
}

type PropertiesResponse struct {
	Results []lite.Property
}

type PropertiesByTransactionType struct {
	Results []lite.ListPropertiesByTransactionTypeRow
}

type IGetPropertiesService interface {
	ListFeaturedProperties(context.Context) (*[]lite.ListFeaturedPropertiesRow, error)
	ListProperties(context.Context) (*[]lite.Property, error)
	ListPropertiesByTransactionType(ctx context.Context, transactionType string) (*[]lite.ListPropertiesByTransactionTypeRow, error)
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

		if _, ok := q[constants.FeaturedQueryKey]; ok {
			if q[constants.FeaturedQueryKey][0] == "true" {
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
		}

		// TODO should you support both featured and propertyType filters?
		if _, ok := q[constants.TransactionTypeQueryKey]; ok {
			if q[constants.TransactionTypeQueryKey][0] != "" {
				var transactionType string = q[constants.TransactionTypeQueryKey][0]
				properties, err := gh.GetPropertiesService.ListPropertiesByTransactionType(r.Context(), transactionType)
				if err != nil {
					gh.Log.Error("error trying to get featured properties", "error", err)
					utils.ReplyError(w, r, http.StatusInternalServerError, fmt.Sprintf("error trying to get the featured properties %v", err))
					return
				}
				response := PropertiesByTransactionType{
					Results: *properties,
				}

				utils.ReplySuccess(w, r, http.StatusOK, response)
				return
			}
		}

		properties, err := gh.GetPropertiesService.ListProperties(r.Context())
		if err != nil {
			gh.Log.Error("error trying to get list of properties", "error", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, fmt.Sprintf("error trying to get list of properties %v", err))
			return
		}
		response := PropertiesResponse{
			Results: *properties,
		}

		utils.ReplySuccess(w, r, http.StatusOK, response)
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "")
}
