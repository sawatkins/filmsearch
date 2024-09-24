window.onload = function () {
    document.getElementById("search-form").addEventListener("submit", handleSearchInput);
    setupExampleSearches();
}

function handleSearchInput(event) {
    event.preventDefault(); 
    const searchQuery = document.getElementById("search-input").value;
    if (searchQuery.trim() === '') {
        return;
    }
    let params = new URLSearchParams();
    params.append("q", searchQuery)
    window.location.href = "/search?" + params
}

function setupExampleSearches() {
    const exampleSearches = document.querySelectorAll('.example-search');
    exampleSearches.forEach(search => {
        search.addEventListener('click', function() {
            const query = this.textContent;
            document.getElementById('search-input').value = query;
            handleSearchInput(new Event('submit'));
        });
    });
}
