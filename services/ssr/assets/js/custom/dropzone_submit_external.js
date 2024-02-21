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
          for (const file of myDropzone.files) {
            let presingResp = await FetchPresign(file.name);
            let respParse = await presingResp.json();
            console.log(file);
            let s3Resp = await fetch(respParse.presignedUrl, {
              method: "PUT",
              headers: {
                "Content-Type": "binary/octet-stream", // Or file.type for the actual MIME type
              },
              body: file,
            });

            CreateHiddenInputWithInfo(
              respParse.keyName,
              imagesFormKey,
              myDropzone.element
            );
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
        $("#myModal .modal-body").html("Success editing the property");
        $("#myModal").modal("show");
        if (response.ok !== undefined && !response.ok) {
          $("#myModal .modal-body").html("Failed to add the property");
          $("#myModal").modal("show");
        }
        // dispatch a event that the page has loaded success
        var submitEvent = new CustomEvent("submitresponse", {
          detail: { key: "loadsuccess" },
        });
        window.dispatchEvent(submitEvent);
      });
      this.on("errormultiple", function (files, response) {
        if (response.ok != undefined && !response.ok) {
          $("#myModal .modal-body").html("Failed editing the property");
          $("#myModal").modal("show");
        }
        if (typeof response === "string") {
          document.documentElement.innerHTML = response;
        }
      });
    },
  };
}
