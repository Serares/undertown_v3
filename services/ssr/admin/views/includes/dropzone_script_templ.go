// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package includes

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"

func DropZone(images []string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_DropZone_7939`,
		Function: `function __templ_DropZone_7939(images){Dropzone.options.uploadForm = {
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
            type: "webp"
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
          myDropzone.processQueue();
        });
      this.on("sendingmultiple", function () {
        console.log("Sending multiple");
      });
      this.on("addedfile", function(file){ 
        console.log(file)
      });
      this.on("successmultiple", function (files, response) {
        if (!response.ok) {
                // console.log("network response", response)
            }
            if (response.redirect) {
                window.location.href = response.url;
            }
            // return response.json();
            document.documentElement.innerHTML = response
            // dispatch a event that the page has loaded success
            var submitEvent = new CustomEvent('submitresponse', {
                    detail: { key: 'loadsuccess' }
                  });
            window.dispatchEvent(submitEvent)
      });
      this.on("errormultiple", function (files, response) {});
    },
  };}`,
		Call:       templ.SafeScript(`__templ_DropZone_7939`, images),
		CallInline: templ.SafeScriptInline(`__templ_DropZone_7939`, images),
	}
}
