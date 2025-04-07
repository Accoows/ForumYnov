// SCRIPT DE GESTION DES LIKES/DISLIKES NON UTILISES

document.addEventListener("DOMContentLoaded", () => {
    function setupLikeSystem(selector, isPost) {
        document.querySelectorAll(selector).forEach(control => {
            const postID = control.dataset.postId || control.dataset.parentPostId;
            const commentID = control.dataset.commentId || null;

            control.querySelectorAll("button").forEach(btn => {
                btn.addEventListener("click", async (e) => {
                    e.preventDefault();
                    const action = btn.dataset.action;

                    const formData = new URLSearchParams();
                    formData.append("action", action);
                    if (isPost) {
                        formData.append("post_id", postID);
                    } else {
                        formData.append("comment_id", commentID);
                        formData.append("parent_post_id", postID);
                    }

                    const res = await fetch("/like", {
                        method: "POST",
                        headers: { "Content-Type": "application/x-www-form-urlencoded" },
                        body: formData
                    });

                    if (res.ok) {
                        // Reload simple pour l'instant
                        location.reload();
                    } else {
                        alert("Erreur lors du traitement du like/dislike");
                    }
                });
            });
        });
    }

    setupLikeSystem(".like-controls-post", true);
    setupLikeSystem(".like-controls-comment", false);
});
