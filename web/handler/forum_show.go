package handler

import "net/http"

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

	h.renderLayout(w, r, "boards", map[string]interface{}{
		"forum":  forum,
		"boards": boards,
	})
}
