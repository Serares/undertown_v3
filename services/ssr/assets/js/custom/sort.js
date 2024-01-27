(function () {
  try {
    let sortFormElement = document.querySelector("#sort_form");
    let selectElement = document.querySelector("#select_sort_type");
    selectElement.addEventListener("change", function () {
      sortFormElement.submit();
    });
  } catch (err) {
    console.log(err);
  }
})();
