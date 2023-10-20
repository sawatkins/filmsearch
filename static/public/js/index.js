window.onload = function () {
    console.log("loaded index.js");
    document.getElementById("search-form").addEventListener("submit", handleSearchInput);
}

function handleSearchInput(event) {
    event.preventDefault(); 
    const searchQuery = document.getElementById("search-input").value;
    console.log("query:" + searchQuery);
    let params = new URLSearchParams();
    params.append("q", searchQuery)
    window.location.href = "/search?" + params

    // fetch("/search?" + params)
    //     .then(response => response.json())
    //     .then(data => {
    //         console.log(data)
    //     })
    //     .catch(error => {
    //         console.log(error)
    //     })
}