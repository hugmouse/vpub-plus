package handler

import (
	"fmt"
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) removePost(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)
	switch r.Method {
	case "GET":
		id := RouteInt64Param(r, "postId")
		post, err := h.storage.PostById(id)
		if err != nil {
			serverError(w, err)
			return
		}
		v := NewView(w, r, "confirm_remove_post")
		v.Set("post", post)
		v.Render()
	case "POST":
		post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
		if err != nil {
			serverError(w, err)
			return
		}

		v := NewView(w, r, "confirm_remove_post")
		v.Set("post", post)

		post.User = user
		if err = h.storage.DeletePost(post); err != nil {
			v.Set("errorMessage", "Unable to delete post: "+err.Error())
			v.Render()
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/topics/%d", post.TopicId), http.StatusFound)
	}
}
