package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/undertown_v3/utils"
)

type AdminLogin struct {
	Log          *slog.Logger
	LoginService service.LoginService
}

func NewLoginHandler(log *slog.Logger, service service.LoginService) *AdminLogin {
	return &AdminLogin{
		Log:          log,
		LoginService: service,
	}
}

func (h *AdminLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// get the email and password
		email := r.FormValue("email")
		password := r.FormValue("password")
		// check if the email and password are valid
		if email == "" || password == "" {
			utils.ReplyError(w, r, http.StatusBadRequest, "invalid email or password")
			return
		}
		token, err := h.LoginService.Login(email, password)
		if err != nil {
			utils.ReplyError(w, r, http.StatusUnauthorized, "invalid email or password")
			return
		}
		cookieExpiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: cookieExpiration,
		}
		http.SetCookie(w, &cookie)
		utils.ReplySuccess(w, r, http.StatusOK, "login successful")
		return
	}
	if r.Method == http.MethodGet {
		viewLogin(w, r)
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
}

func viewLogin(w http.ResponseWriter, r *http.Request) {
	views.Login().Render(r.Context(), w)
}
