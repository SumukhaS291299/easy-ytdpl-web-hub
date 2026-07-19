document.addEventListener("DOMContentLoaded", function () {
    const forms = document.querySelectorAll("form");
    forms.forEach(form => {
        form.addEventListener("submit", function () {
            const overlay = document.getElementById("loading-overlay");
            const text = document.getElementById("loading-text");
            if (overlay && text) {
                if (form.classList.contains("probe-form")) {
                    text.innerText = "Probing media streams...";
                } else if (form.classList.contains("download-form")) {
                    text.innerText = "Downloading and processing media...";
                } else {
                    text.innerText = "Processing request...";
                }
                overlay.classList.add("active");
            }
        });
    });
});
