package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils"
)

type AdminLogin struct {
	Log          *slog.Logger
	LoginService *service.LoginService
}

func NewLoginHandler(log *slog.Logger, service *service.LoginService) *AdminLogin {
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
			viewLogin(w, r, types.LoginProps{
				ErrorMessage: "Email or password cannot be null",
			})
			return
		}
		token, err := h.LoginService.Login(email, password)
		if err != nil {
			h.Log.Error("error invalid response", "err", err)
			viewLogin(w, r, types.LoginProps{
				ErrorMessage: "Invalid email or password",
			})
			return
		}
		cookieExpiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: cookieExpiration,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, types.ListPath, http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		// TODO if it's a redirect
		// you'll have to use cookies to send a message to login

		viewLogin(w, r, types.LoginProps{})
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
}

func viewLogin(w http.ResponseWriter, r *http.Request, props types.LoginProps) {
	views.Login(
		types.BasicIncludes{
			Header:        components.Header("Login"),
			BannerSection: components.BannerSection(includesTypes.BannerSectionProps{Title: "Login"}),
			Preload:       components.Preload(),
			Navbar:        components.Navbar(includesTypes.NavbarProps{Path: "/login", IsAdmin: true}),
			Footer:        components.Footer(),
			Scripts:       components.Scripts(),
		},
		props,
	).Render(r.Context(), w)
}
