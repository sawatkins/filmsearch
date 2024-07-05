window.onload = function () {
    console.log("loaded script.js");
    document.querySelectorAll(".tmdb_info").forEach(element => {
        if (element.href === "not found") {
            element.textContent = "no external info found";
            element.removeAttribute("href");
        }
    });
}
