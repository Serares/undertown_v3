package handler

import (
	"fmt"
	"io"
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
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		var err error
		var isLocal = os.Getenv("IS_LOCAL")
		v1Request, ok := algnhsa.APIGatewayV1RequestFromContext(r.Context())
		if !ok {
			h.Log.Error("Error trying to get the APIV1Request Context", "error", err)
		}
		// h.Log.Info("DEBUGGING", "request", v1Request)
		// get the user id from the apiv1 request context
		userId, ok := v1Request.RequestContext.Authorizer["userId"].(string)
		if !ok && isLocal != "true" {
			h.Log.Error("the request is missing the userId from the authorizer")
			utils.ReplyError(w, r, http.StatusForbidden, "Can't find the user id in the request")
			return
		}

		// ❗for local testing purposes
		if userId == "" && isLocal == "true" {
			userId = "c8fd42e9-7c8f-4bf0-b818-f6bb96304e92"
		}

		const maxFileUpload = 10 << 20
		var imagesNames []string
		h.Log.Info("Content-Type", "value", r.Header.Get("Content-Type"))
		if err := r.ParseMultipartForm(maxFileUpload); err != nil {
			h.Log.Error("Error parsing the form", "err", err)
			body, err := io.ReadAll(r.Body)
			if err != nil {
				h.Log.Error("error reading the body", "err", err)
			}
			h.Log.Info("DEBUG THE BODY", "body", string(body))
			utils.ReplyError(w, r, http.StatusExpectationFailed, "files are too large")
			return
		}
		images := r.MultipartForm.File[utils.ImagesFormKey]
		imagesToBeRemoved := r.MultipartForm.Value[utils.DeleteImagesFormKey]

		if len(images) > 0 {
			var processImagesErr error
			if isLocal == "true" {
				imagesNames, processImagesErr = h.SubmitService.ProcessImagesLocal(r.Context(), images)
			} else {
				imagesNames, processImagesErr = h.SubmitService.ProcessImagesS3(r.Context(), images)
			}
			if processImagesErr != nil {
				h.Log.Error("error on processing the images", "error", processImagesErr)
				utils.ReplyError(w, r, http.StatusInternalServerError, "error processing the images")
				return
			}
		}

		if len(imagesToBeRemoved) > 0 {
			var deleteImagesError error
			if isLocal == "true" {
				deleteImagesError = h.SubmitService.DeleteImagesLocal(imagesToBeRemoved)
			} else {
				deleteImagesError = h.SubmitService.DeleteImagesS3(r.Context(), imagesToBeRemoved)
			}
			if deleteImagesError != nil {
				h.Log.Error("error trying to delete the images", "error", err, "images names:", imagesToBeRemoved)
				utils.ReplyError(w, r, http.StatusInternalServerError, "error deleting the image")
				return
			}
		}

		// TODO ❗
		// this pattern of using conditions seems a bit odd
		// check if it's an edit request
		// TODO you will also have a PUT request handled here
		// TODO what if the data process/store failes, the images will be persisted without a property
		// you will have to process and store the property data first? and then upload the images?
		// but you won't have the images paths (get the images names from the multipartForm.File? and store the paths before uploading the images?)
		// think how to solve this
		q := r.URL.Query()
		if _, ok := q[utils.HumanReadableIdQueryKey]; ok {
			humanReadableId := q[utils.HumanReadableIdQueryKey][0]
			err = h.SubmitService.ProcessPropertyUpdateData(r.Context(), imagesNames, imagesToBeRemoved, r.MultipartForm, humanReadableId)
		} else {
			_, _, err = h.SubmitService.ProcessPropertyData(r.Context(), imagesNames, r.MultipartForm, userId)
		}

		fmt.Printf("Uploaded File:")
		if err != nil {
			h.Log.Error("error on processing the property data", "error", err)
			utils.ReplyError(w, r, http.StatusInternalServerError, "error processing the property data")
			return
		}
		utils.ReplySuccess(w, r, http.StatusOK, "Added proeprty succesfully")
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
}
