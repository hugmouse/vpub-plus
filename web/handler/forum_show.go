package handler

import (
	"net/http"
)

func (h *Handler) showForumView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumByID(id)
	if err != nil {
		notFound(w)
		return
	}

	boards, err := h.storage.BoardsByForumID(forum.ID)
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
