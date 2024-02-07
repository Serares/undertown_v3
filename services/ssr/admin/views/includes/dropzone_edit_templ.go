// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package includes

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"

func DropzoneEdit(images []string, url string, deleteImagesFormKey string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_DropzoneEdit_5f78`,
		Function: `function __templ_DropzoneEdit_5f78(images, url, deleteImagesFormKey){Dropzone.options.uploadForm = {
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
      if (images && images.length) {
      images.forEach((path)=>{
        // those are just some mock data
        const fileData = {
            name: path,
            size: 12345,
            type: "webp",
            status: "s3"
        }
        myDropzone.displayExistingFile( fileData, path)
      })
      }
      this.element
        .querySelector("button[type=submit]")
        .addEventListener("click", function (e) {
          e.preventDefault();
          e.stopPropagation();
          document.querySelector(".dz-hidden-input").setAttribute("name", "images")
          // Have to handle the case where the user is not adding any new images when modifying the property
          if (myDropzone.getQueuedFiles().length > 0) {
            myDropzone.processQueue();
          } else {
            // manually send the AJAX request
            let formData = new FormData(myDropzone.element)
            fetch(url, {
            method: "POST",
            body: formData,
            }
            )
            .then(response => {
              if (response.ok !== undefined && !response.ok) {
                myDropzone.emit("errormultiple", [], response)
                return
              }
              myDropzone.emit("successmultiple", [], response)
            })
            .catch(err=> {
              console.log("error sending thew AJAX request")
            })
          }
        });
      this.on("sendingmultiple", function () {
        console.log("Sending multiple");
      });
      this.on("addedfile", function(file){ 
        console.log(file)
      });
      this.on("successmultiple", function (files, response, xhr) {
        
        if (response.redirected) {
            window.location.href = response.url;
            return;
          }
          // return response.json();
          document.documentElement.innerHTML = response
          if (response.ok != undefined && response.ok) {
                $("#myModal .modal-body").html("Success editing the property")
                $("#myModal").modal("show")
          }
          // dispatch a event that the page has loaded success
          var submitEvent = new CustomEvent('submitresponse', {
                  detail: { key: 'loadsuccess' }
                });
          window.dispatchEvent(submitEvent)
      });
      this.on("errormultiple", function (files, response, xhr) {
        if (response.ok != undefined && !response.ok) {
                $("#myModal .modal-body").html("Failed editing the property")
                $("#myModal").modal("show")
          }
      });
      this.on("removedfile", function(file) {
        // check if status of the image is S3 and add it to the remove list
        if (file.status && file.status.toLowerCase() === "s3") {
          // create inputs with the same name for each delete file
          // because inputs with the same names will be added as an array of values with the same key
          let newInput = document.createElement("input");
          newInput.type = "text"
          newInput.id = "input-" + Math.random().toString(36).substr(2,9)
          newInput.setAttribute("name", deleteImagesFormKey);
          newInput.setAttribute("style", "visibility: hidden;");
          newInput.classList.add("hidden");
          newInput.value = file.name;
          this.element.appendChild(newInput);
        }
      })
    },
  };}`,
		Call:       templ.SafeScript(`__templ_DropzoneEdit_5f78`, images, url, deleteImagesFormKey),
		CallInline: templ.SafeScriptInline(`__templ_DropzoneEdit_5f78`, images, url, deleteImagesFormKey),
	}
}
