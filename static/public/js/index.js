window.onload = function () {
    console.log("loaded index.js");
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
    const searchQuery = document.getElementById("search-input").value;
    console.log("query:" + searchQuery);
    let params = new URLSearchParams();
    params.append("q", searchQuery)
    window.location.href = "/search?" + params
}
