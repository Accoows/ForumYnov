// Remet la page à la position précédente après rechargement
window.addEventListener("load", () => {
    const savedY = sessionStorage.getItem("scrollY");
    if (savedY !== null) {
        window.scrollTo(0, parseInt(savedY));
        sessionStorage.removeItem("scrollY");
    }
});

// Sauvegarde la position du scroll avant d’envoyer un formulaire
document.querySelectorAll("form").forEach(form => {
    form.addEventListener("submit", () => {
        sessionStorage.setItem("scrollY", window.scrollY);
    });
});
