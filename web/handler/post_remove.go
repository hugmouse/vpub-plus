package handler

import (
	"fmt"
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) removePost(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

	post, err := h.storage.PostByID(RouteInt64Param(r, "postId"))
	if err != nil {
		serverError(w, err)
		return
	}

	topic, err := h.storage.TopicByID(post.TopicID)
	if err != nil {
		serverError(w, err)
		return
	}

	board, err := h.storage.BoardByID(topic.BoardID)
	if err != nil {
		serverError(w, err)
		return
	}

	if !canAccessForum(board.Forum, user) {
		forbidden(w)
		return
	}

	v := NewView(w, r, "confirm_remove_post")
	v.Set("post", post)

	switch r.Method {
	case "GET":
		v.Render()
	case "POST":
		post.User = user
		if err = h.storage.DeletePost(post); err != nil {
			v.Set("errorMessage", "Unable to delete post: "+err.Error())
			v.Render()
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", post.TopicID), http.StatusFound)
	}
}
