package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Serares/api/register/service"
	"github.com/Serares/api/register/types"
)

const (
	ErrorUnmarshal = "error trying to unmarshal the request"
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
			rh.Log.Error("%s - %w", ErrorUnmarshal, err)
			// TODO handle response error
		}
		err = rh.RegisterService.PersistUser(r.Context(), &user)
		if err != nil {
			// TODO handle error response
			rh.Log.Error("error persisting user", "type", types.ErrorPersistingUser, "error", err)
		}
		// TODO success response
	}
}
