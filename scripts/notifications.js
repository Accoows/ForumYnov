window.addEventListener("DOMContentLoaded", () => {
  const toast = document.getElementById("notification-toast");
  const messageEl = document.getElementById("notification-message");
  const closeBtn = document.getElementById("close-notification-toast");  

  const getCookie = (name) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
  };

  const msg = getCookie("notif_msg");
  const type = getCookie("notif_type"); // success / info / error

  if (msg) {
    toast.classList.remove("hidden");
    messageEl.textContent = decodeURIComponent(msg.replace(/\+/g, ' '));

    toast.classList.remove("bg-green", "bg-blue", "bg-red");
    if (type === "success") toast.classList.add("bg-green");
    else if (type === "info") toast.classList.add("bg-blue");
    else toast.classList.add("bg-red"); // Default : error

    toast.style.animation = "slideIn 0.4s ease";

    // Delete cookies
    document.cookie = "notif_msg=; Max-Age=0; path=/;";
    document.cookie = "notif_type=; Max-Age=0; path=/;";

    // Automatic disappearance
    setTimeout(() => {
      toast.style.animation = "slideOut 0.4s ease forwards";
      setTimeout(() => {
        toast.classList.add("hidden");
        toast.style.animation = "";
      }, 400);
    }, 5000);

    // Manual closing
    closeBtn.addEventListener("click", () => {
      toast.style.animation = "slideOut 0.4s ease forwards";
      setTimeout(() => {
        toast.classList.add("hidden");
        toast.style.animation = "";
      }, 400);
    });
  }
});
