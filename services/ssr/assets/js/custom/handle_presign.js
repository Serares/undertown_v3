// transactionType in case of submit
function FetchPresign(fileName) {
  return fetch("/presign", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      fileName,
    }),
  }).then((response) => response.json());
}

//
function uploadFileToS3(presignedUrl, file) {
  return fetch(presignedUrl, {
    method: "PUT",
    headers: {
      "Content-Type": "binary/octet-stream", // Adjust this according to your file type
    },
    body: file,
  });
}

// imageName: string
// imagesFormKey: string
// formElement: domElement
function CreateHiddenInputWithInfo(value, inputName, attachToElement) {
  // create inputs with the same name for each delete file
  // because inputs with the same names will be added as an array of values with the same key
  let newInput = document.createElement("input");
  newInput.type = "text";
  newInput.id = "input-" + Math.random().toString(36).substr(2, 9);
  newInput.setAttribute("name", inputName);
  newInput.setAttribute("style", "visibility: hidden;");
  newInput.classList.add("hidden");
  newInput.value = value;
  attachToElement.appendChild(newInput);
}
