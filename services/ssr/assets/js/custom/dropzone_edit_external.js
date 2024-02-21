function InitEditDropzone(
  imagesPaths,
  url,
  deleteImagesFormKey,
  imagesFormKey
) {
  Dropzone.options.uploadForm = {
    autoProcessQueue: false,
    addRemoveLinks: true,
    uploadMultiple: true,
    parallelUploads: 100,
    maxFiles: 100,
    previewsContainer: ".dropzone-previews",
    hiddenInputContainer: ".browse_submit",
    paramName: "#images-input",
    init: function () {
      var myDropzone = this;
      if (imagesPaths && imagesPaths.length) {
        imagesPaths.forEach((path) => {
          // those are just some mock data
          imagePathSplitted = path.split("/");
          const fileData = {
            name: imagePathSplitted[imagePathSplitted.length - 1],
            size: 12345,
            type: "webp",
            status: "s3",
          };
          myDropzone.displayExistingFile(fileData, path);
        });
      }
      this.element
        .querySelector("button[type=submit]")
        .addEventListener("click", async function (e) {
          e.preventDefault();
          e.stopPropagation();
          document
            .querySelector(".dz-hidden-input")
            .setAttribute("name", "thumbnails");
          // Have to handle the case where the user is not adding any new images when modifying the property
          // manually send the AJAX request
          let formData = new FormData(myDropzone.element);
          const presignedUrlsPromises = myDropzone.files.map((file) =>
            FetchPresign(file.name)
          );
          try {
            const presignedUrls = await Promise.all(presignedUrlsPromises);
            const uploadPromises = presignedUrls.map((urlData, index) => {
              const { presignedUrl, keyName } = urlData; // Assuming the response includes the presigned URL
              CreateHiddenInputWithInfo(
                keyName,
                imagesFormKey,
                myDropzone.element
              );
              return uploadFileToS3(presignedUrl, myDropzone.files[index]);
            });
            await Promise.all(uploadPromises);
            console.log("All files uploaded successfully.");
          } catch (err) {
            myDropzone.emit("errormultiple", [], err);
            return;
          }
          formData = new FormData(myDropzone.element);
          //  delete anything that's a file from the form
          for (const [key, value] of formData.entries()) {
            if (value instanceof File) {
              formData.delete(key);
            }
          }
          fetch(url, {
            method: "POST",
            body: formData,
          })
            .then((response) => {
              if (response.ok !== undefined && !response.ok) {
                myDropzone.emit("errormultiple", [], response);
                return;
              }
              return response.text();
            })
            .then((data) => {
              myDropzone.emit("successmultiple", [], data);
            })
            .catch((err) => {
              console.log("error sending thew AJAX request");
            });
        });
      this.on("sendingmultiple", function () {});
      this.on("addedfile", function (file) {});
      this.on("successmultiple", function (files, response, xhr) {
        if (response.redirected) {
          window.location.href = response.url;
          return;
        }
        // return response.json();
        document.documentElement.innerHTML = response;
        $("#myModal .modal-body").html("Success editing the property").css({
          color: "green",
          "font-weight": "900",
        });
        $("#myModal").modal("show");
        // dispatch a event that the page has loaded success
        var submitEvent = new CustomEvent("submitresponse", {
          detail: { key: "loadsuccess" },
        });
        window.dispatchEvent(submitEvent);
      });
      this.on("errormultiple", function (files, response, xhr) {
        $("#myModal .modal-body")
          .html(
            "Failed editing the property, check the console for debug clues"
          )
          .css({
            color: "red",
            "font-weight": "900",
          });
        $("#myModal").modal("show");
        console.log(response);
      });
      this.on("removedfile", function (file) {
        // check if status of the image is S3 and add it to the remove list
        if (file.status && file.status.toLowerCase() === "s3") {
          CreateHiddenInputWithInfo(
            file.name,
            deleteImagesFormKey,
            this.element
          );
        }
      });
    },
  };
}
