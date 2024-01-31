package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	includesTypes "github.com/Serares/undertown_v3/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils"
)

type AdminSubmit struct {
	Log           *slog.Logger
	SubmitService *service.SubmitService
}

func NewSubmitHandler(log *slog.Logger, submitService *service.SubmitService) *AdminSubmit {
	return &AdminSubmit{
		Log:           log,
		SubmitService: submitService,
	}
}

func (h *AdminSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			h.Log.Info("Content type", "type", r.Header.Get("Content-Type"))
			utils.ReplyError(w, r, http.StatusBadRequest, "Invalid content type")
			return
		}
		var errorMessage string
		var successMessage string

		cookie, err := r.Cookie("token")
		if err != nil {
			errorMessage = "Invalid authentication"
			viewSubmit(w, r, types.SubmitProps{Message: successMessage, ErrorMessage: errorMessage})
			return
		}

		err = h.SubmitService.Submit(r, cookie.Value, "", false)

		if err != nil {
			h.Log.Error("failed to submit", err)
			errorMessage = "Failed to submit the property"
			viewSubmit(w, r, types.SubmitProps{Message: successMessage, ErrorMessage: errorMessage})
			return
		}
		if errorMessage == "" {
			successMessage = "Success submitting the property"
		}
		viewSubmit(w, r, types.SubmitProps{Message: successMessage, ErrorMessage: errorMessage})
		return
	}
	if r.Method == http.MethodGet {
		viewSubmit(w, r, types.SubmitProps{Message: "Property get success", ErrorMessage: ""})
		return
	}
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
}

func viewSubmit(w http.ResponseWriter, r *http.Request, submitProps types.SubmitProps) {
	views.Submit(types.BasicIncludes{
		Header: components.Header("Submit"),
		BannerSection: components.BannerSection(includesTypes.BannerSectionProps{
			Title: "Submit",
		},
		),
		Preload: components.Preload(),
		Navbar: components.Navbar(includesTypes.NavbarProps{
			Path:    "/",
			IsAdmin: true,
		}),
		Footer:  components.Footer(),
		Scripts: components.Scripts(),
	},
		submitProps,
	).Render(r.Context(), w)
}
