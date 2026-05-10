package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/request"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	user := request.GetUserContextKey(r)
	forums := forumFromBoards(boards)
	var visible []model.Forum
	for _, f := range forums {
		if canSeeForum(f, user) {
			visible = append(visible, f)
		}
	}

	v := NewView(w, r, "index")
	v.Set("forums", visible)
	v.Render()
}
