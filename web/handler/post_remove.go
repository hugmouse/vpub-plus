package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
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
		h.renderLayout(w, r, "confirm_remove_post", map[string]interface{}{
			"post":           post,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	case "POST":
		post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
		if err != nil {
			serverError(w, err)
			return
		}
		post.User = user
		err = h.storage.DeletePost(post)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", post.TopicId), http.StatusFound)
	}
}
