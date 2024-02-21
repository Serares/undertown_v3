function InitSubmitDropzone(imagesFormKey) {
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
      this.element
        .querySelector("button[type=submit]")
        .addEventListener("click", async function (e) {
          e.preventDefault();
          e.stopPropagation();
          document
            .querySelector(".dz-hidden-input")
            .setAttribute("name", "thumbnails");
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

          // initialize the formdata again with the images names
          formData = new FormData(myDropzone.element);

          //  delete anything that's a file from the form
          for (const [key, value] of formData.entries()) {
            if (value instanceof File) {
              formData.delete(key);
            }
          }
          // send the rest of the form inputs
          fetch("/submit", {
            method: "POST",
            body: formData,
          })
            .then((resp) => {
              return resp.text();
            })
            .then((data) => {
              myDropzone.emit("successmultiple", [], data);
            })
            .catch((err) => {
              console.log("error sending thew AJAX request");
            });
        });
      this.on("sendingmultiple", function () {
        console.log("Sending multiple");
      });
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
        if (response.ok !== undefined && !response.ok) {
          $("#myModal .modal-body").html("Failed to add the property").css({
            color: "red",
            "font-weight": "900",
          });
          $("#myModal").modal("show");
        }
        // dispatch a event that the page has loaded success
        var submitEvent = new CustomEvent("submitresponse", {
          detail: { key: "loadsuccess" },
        });
        window.dispatchEvent(submitEvent);
      });
      this.on("errormultiple", function (files, response) {
        $("#myModal .modal-body")
          .html(
            "Failed adding the property, try again later or check the console for debugg clues"
          )
          .css({
            color: "green",
            "font-weight": "900",
          });
        $("#myModal").modal("show");
        console.log(response);

        if (typeof response === "string") {
          document.documentElement.innerHTML = response;
        }
      });
    },
  };
}
