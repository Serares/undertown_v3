package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/services/api/addProperty/service"
	"github.com/Serares/undertown_v3/services/api/addProperty/types"
	"github.com/Serares/undertown_v3/utils"
)

type AddPropertyHandler struct {
	Log           *slog.Logger
	SubmitService service.Submit
}

func New(log *slog.Logger, ss service.Submit) *AddPropertyHandler {
	return &AddPropertyHandler{
		Log:           log,
		SubmitService: ss,
	}
}

func (h *AddPropertyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var property types.POSTProperty

		if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
			message := fmt.Sprintf("Invalid JSON: %v", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, message)
			return
		}
		id, hrID, err := h.SubmitService.ProcessProperty(r.Context(), &property)
		if err != nil {
			h.Log.Error("Error processing the property", "url", r.URL, "method", r.Method, "status:", http.StatusInternalServerError, "error:", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, fmt.Sprintf("error trying to persist the order with error: %v", err))
			return
		}
		successReply := types.POSTSuccessResponse{
			PropertyId:      id,
			HumanReadableId: hrID,
		}
		err = utils.ReplySuccess(w, r, http.StatusCreated, successReply)
		if err != nil {
			h.Log.Error("error trying to reply to the request", "error", err)
		}
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, types.ErrorMethodNotSupported)
	return
}
