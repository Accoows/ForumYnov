window.addEventListener("DOMContentLoaded", () => {
  // Get all required DOM elements
  const toast = document.getElementById("notification-toast");
  const messageEl = document.getElementById("notification-message");
  const closeBtn = document.getElementById("close-notification-toast");  

  // Function to read a specific cookie by name
  const getCookie = (name) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
  };

  // Get notification message and type from cookies
  const msg = getCookie("notif_msg");
  const type = getCookie("notif_type"); // success / info / error

  if (msg) {
    // Display the toast notification
    toast.classList.remove("hidden");
    messageEl.textContent = decodeURIComponent(msg.replace(/\+/g, ' '));

    // Remove old color classes and apply new one based on type
    toast.classList.remove("bg-green", "bg-blue", "bg-red");
    if (type === "success") toast.classList.add("bg-green");
    else if (type === "info") toast.classList.add("bg-blue");
    else toast.classList.add("bg-red"); // Default : error

    toast.style.animation = "slideIn 0.4s ease"; // Trigger slide-in animation

    // Delete cookies so they don't persist
    document.cookie = "notif_msg=; Max-Age=0; path=/;";
    document.cookie = "notif_type=; Max-Age=0; path=/;";

    // Automatically hide the toast after 5 seconds
    setTimeout(() => {
      toast.style.animation = "slideOut 0.4s ease forwards";
      setTimeout(() => {
        toast.classList.add("hidden");
        toast.style.animation = "";
      }, 400);
    }, 5000);

    // Allow user to manually close the toast
    closeBtn.addEventListener("click", () => {
      toast.style.animation = "slideOut 0.4s ease forwards";
      setTimeout(() => {
        toast.classList.add("hidden");
        toast.style.animation = "";
      }, 400);
    });
  }
});