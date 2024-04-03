setTimeout(() => {

  let fuzzyInput = document.getElementById("fuzzyInput");
  fuzzyInput.addEventListener("input", (event) => {
    browser.history.search({
      text: "mozilla",
    }).then((myHistory) => {
      historyList = document.getElementById("historyList");
      historyList.innerHTML = JSON.stringify(myHistory, null, 2);
    })
  });

}, 1000);

