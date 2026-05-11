package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) showForumView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}

	user := request.GetUserContextKey(r)
	if !canSeeForum(forum, user) {
		notFound(w)
		return
	}
	if !canAccessForum(forum, user) {
		forbidden(w)
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
