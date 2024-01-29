package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/ssr/includes/components"
	includesTypes "github.com/Serares/ssr/includes/types"
	"github.com/Serares/undertown_v3/utils"
)

type AdminSubmit struct {
	Log           *slog.Logger
	SubmitService service.SubmitService
}

func NewSubmitHandler(log *slog.Logger, submitService service.SubmitService) *AdminSubmit {
	return &AdminSubmit{
		Log:           log,
		SubmitService: submitService,
	}
}

func (h *AdminSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Got the request")
	if r.Method == http.MethodPost {
		h.Log.Info("Got in the POST request", "req", r.Header.Get("Content-Type"))
		if r.Header.Get("Content-Type") != "multipart/form-data" {
			h.Log.Info("Content type", "type", r.Header.Get("Content-Type"))
			utils.ReplyError(w, r, http.StatusBadRequest, "Invalid content type")
			return
		}
		// 32 mb is a lot
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.ReplyError(w, r, http.StatusBadRequest, "Invalid form data")
			return
		}
		// check if property exists
		property := r.FormValue("property")
		if property == "" {
			utils.ReplyError(w, r, http.StatusBadRequest, "Property is required")
			return
		}
		_, _, err := r.FormFile("images")
		if err != nil {
			utils.ReplyError(w, r, http.StatusBadRequest, "Images are required")
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.ReplyError(w, r, http.StatusBadRequest, "Invalid token")
			return
		}
		err = h.SubmitService.Submit(r.Body, cookie.Value)
		if err != nil {
			h.Log.Error("failed to submit", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, "Failed to submit")
			return
		}

		viewSubmit(w, r, types.SubmitProps{Message: "Property submitted success", ErrorMessage: ""})
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
			Path: "/",
		}),
		Footer:  components.Footer(),
		Scripts: components.Scripts(),
	},
		submitProps,
	).Render(r.Context(), w)
}
