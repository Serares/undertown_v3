package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Serares/undertown_v3/services/api/login/service"
	"github.com/Serares/undertown_v3/utils"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	JWTToken string `json:"jwtToken,omitempty"` // this is returned in base64
	Error    string `json:"error,omitempty"`
}

type LoginHandler struct {
	Log          *slog.Logger
	LoginService *service.LoginService
}

func NewLoginHandler(log *slog.Logger, ls *service.LoginService) *LoginHandler {
	return &LoginHandler{
		Log:          log,
		LoginService: ls,
	}
}

func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var userInfo Request
		err := json.NewDecoder(r.Body).Decode(&userInfo)
		if err != nil {
			utils.ReplyError(w, r, http.StatusInternalServerError, "error decoding the request")
			return
		}
		jwt, err := lh.LoginService.LoginUser(r.Context(), userInfo.Email, userInfo.Password)
		if err != nil {
			lh.Log.Error("error trying to login the user", "error", err)
			utils.ReplyError(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		utils.ReplySuccess(w, r, http.StatusAccepted, Response{JWTToken: jwt})
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "method not allowed")
}
