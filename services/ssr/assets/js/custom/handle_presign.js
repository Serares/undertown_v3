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

// Function to chunk an array into batches
function chunkArray(array, size) {
  const chunkedArr = [];
  for (let i = 0; i < array.length; i += size) {
    chunkedArr.push(array.slice(i, i + size));
  }
  return chunkedArr;
}

function processPresignBatch(filesBatch) {
  const batchPromises = filesBatch.map((file) => {
    return FetchPresign(file.name).then((presignedData) => ({
      //
      presignedData,
      file,
    }));
  });
  return Promise.all(batchPromises);
}

function processS3UploadBatch(presignDataFilesBatch) {
  const batchPromises = presignDataFilesBatch.map(({ presignedData, file }) => {
    let { presignedUrl, keyName } = presignedData;
    return uploadFileToS3(presignedUrl, file).then(() => keyName);
  });
  return Promise.all(batchPromises);
}

// Function to process requests in batches
function ProcessPresignInBatches(files, batchSize) {
  const fileBatches = chunkArray(files, batchSize);
  let promiseChain = Promise.resolve(); // Start with a resolved promise for chaining
  const results = []; // Array to collect all results
  const keyNames = [];

  fileBatches.forEach((fileBatch) => {
    promiseChain = promiseChain
      .then(() => processPresignBatch(fileBatch))
      .then((presignDataFilesResults) => {
        console.log("Batch results:", presignDataFilesResults);
        results.push(...presignDataFilesResults); // Collect results
        // Optional: return a new promise if you want to wait between batches
        return presignDataFilesResults;
      })
      .then((results) => {
        return processS3UploadBatch(results);
      })
      .then((keyNamesResults) => {
        keyNames.push(...keyNamesResults);
        return new Promise((resolve) => setTimeout(resolve, 500));
      })
      .catch((err) => {
        return err;
      });
  });

  return promiseChain.then(() => keyNames);
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
