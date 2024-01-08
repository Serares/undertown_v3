package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/services/api/register/service"
	"github.com/Serares/undertown_v3/services/api/register/types"
	"github.com/Serares/undertown_v3/utils"
)

type RegisterHandler struct {
	Log             *slog.Logger
	RegisterService service.UserService
}

func NewRegisterHandler(log *slog.Logger, registerService service.UserService) *RegisterHandler {
	return &RegisterHandler{
		Log:             log,
		RegisterService: registerService,
	}
}

func (rh *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// unmarshal the request
		var user types.PostUserRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			rh.Log.Error("error trying to unmarshal the request", "error:", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, "error on POST request")
			return
		}
		err = rh.RegisterService.PersistUser(r.Context(), &user)
		if err != nil {
			rh.Log.Error("error persisting user", "type", types.ErrorPersistingUser, "error", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		utils.ReplySuccess(w, r, http.StatusAccepted, "success persisting user")
	}
}
