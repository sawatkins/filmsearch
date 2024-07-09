document.addEventListener('htmx:afterOnLoad', function (event) {
    document.querySelectorAll(".tmdb_info").forEach(element => {
        if (element.href.includes("notfound")) {
            element.textContent = "no external info found";
            element.removeAttribute("href");
        }
    });
});