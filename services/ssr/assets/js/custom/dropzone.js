(function () {
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
        .addEventListener("click", function (e) {
          // e.preventDefault();
          e.stopPropagation();
          myDropzone.processQueue();
        });
      this.on("sendingmultiple", function () {
        console.log("Sending multiple");
      });
      this.on("successmultiple", function (files, response) {});
      this.on("errormultiple", function (files, response) {});
    },
  };
})();
