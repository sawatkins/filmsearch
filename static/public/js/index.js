window.onload = function () {
    console.log("loaded index.js");
    hideSpinner();
    document.getElementById("search-form").addEventListener("submit", handleSearchInput);

    document.querySelectorAll('.example-searches a').forEach(link => {
        link.addEventListener('click', function(event) {
            event.preventDefault();
            document.getElementById("search-input").value = this.textContent;
            handleSearchInput(event);
        });
    });
}

function handleSearchInput(event) {
    event.preventDefault(); 
    displaySpinner();
    const searchQuery = document.getElementById("search-input").value;
    console.log("query:" + searchQuery);
    let params = new URLSearchParams();
    params.append("q", searchQuery)
    window.location.href = "/search?" + params
}

function displaySpinner() {
    const spinner = document.getElementById('spinner');
    spinner.style.display = 'block';
}

function hideSpinner() {
    const spinner = document.getElementById('spinner');
    spinner.style.display = 'none';
}
