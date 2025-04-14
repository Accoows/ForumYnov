//Dark Mode Toggle
document.getElementById('dark-mode-toggle').addEventListener('click', function() {
    document.body.classList.toggle('dark-mode');
});

function confirmDelete() {
    return confirm("Are you sure you want to delete this post?");
}