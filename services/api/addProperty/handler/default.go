package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Serares/undertown_v3/services/api/addProperty/service"
	"github.com/Serares/undertown_v3/services/api/addProperty/types"
	"github.com/Serares/undertown_v3/utils"
	"github.com/akrylysov/algnhsa"
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
		var err error
		v1Request, ok := algnhsa.APIGatewayV1RequestFromContext(r.Context())
		if !ok {
			h.Log.Error("Error trying to get the APIV1Request Context", "error", err)
		}
		// get the user id from the apiv1 request context
		userId, ok := v1Request.RequestContext.Authorizer["userId"].(string)
		if !ok {
			h.Log.Error("the request is missing the userId from the authorizer")
			utils.ReplyError(w, r, http.StatusForbidden, "Can't find the user id in the request")
			return
		}
		const maxFileUpload = 10 << 20
		var isLocal = os.Getenv("IS_LOCAL")
		var imagesPaths []string
		if err := r.ParseMultipartForm(maxFileUpload); err != nil {
			utils.ReplyError(w, r, http.StatusExpectationFailed, "files are too large")
			return
		}

		files := r.MultipartForm.File["images"]
		if err != nil {
			utils.ReplyError(w, r, http.StatusInternalServerError, "error uploading file to s3")
		}
		if isLocal == "true" {
			imagesPaths, err = h.SubmitService.ProcessPropertyImagesLocal(r.Context(), files)
		} else {
			imagesPaths, err = h.SubmitService.ProcessPropertyImages(r.Context(), files)
		}
		if err != nil {
			h.Log.Error("error on processing the images", "error", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, "error processing the images")
		}

		h.SubmitService.ProcessPropertyData(r.Context(), imagesPaths, r.MultipartForm, userId)
		fmt.Printf("Uploaded File:")

		return
	}
	// successReply := types.POSTSuccessResponse{
	// 	PropertyId:      id,
	// 	HumanReadableId: hrID,
	// }
	// err = utils.ReplySuccess(w, r, http.StatusCreated, successReply)
	// if err != nil {
	// 	h.Log.Error("error trying to reply to the request", "error", err)
	// }
	// return
	utils.ReplyError(w, r, http.StatusMethodNotAllowed, types.ErrorMethodNotSupported)
	return
}
