// Restores the page to the previous position after reloading
window.addEventListener("load", () => {
    const savedY = sessionStorage.getItem("scrollY");
    if (savedY !== null) {
        window.scrollTo(0, parseInt(savedY));
        sessionStorage.removeItem("scrollY");
    }
});

// Saves the scroll position before submitting a form
document.querySelectorAll("form").forEach(form => {
    form.addEventListener("submit", () => {
        sessionStorage.setItem("scrollY", window.scrollY);
    });
});
