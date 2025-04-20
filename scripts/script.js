//Dark Mode Toggle
document.getElementById('dark-mode-toggle').addEventListener('click', function() {
    document.body.classList.toggle('dark-mode');
});

function confirmDelete() { // Confirm deletion of a post in view_post.html
    return confirm("Are you sure you want to delete this post?");
}

function enableEdit(input) {
    console.log("enableEdit triggered on", input);
    input.removeAttribute("readonly");
    input.focus();
}