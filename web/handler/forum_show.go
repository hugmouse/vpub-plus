package handler

import (
	"net/http"
)

func (h *Handler) showForumView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumById(id)
	if err != nil {
		notFound(w)
		return
	}

	boards, err := h.storage.BoardsByForumId(forum.Id)
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "boards")
	v.Set("navigation", navigation{
		Forum: forum,
	})
	v.Set("forum", forum)
	v.Set("boards", boards)
	v.Render()
}
