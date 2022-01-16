package handler

import (
	"html/template"
	"net/http"
	"vpub/model"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)
	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	var forums []model.Forum
	var forum model.Forum
	for i, board := range boards {
		if i == 0 {
			forum.Name = board.Forum.Name
		} else if board.Forum.Name != forum.Name {
			forums = append(forums, forum)
			forum = model.Forum{Name: board.Forum.Name}
		}
		forum.Boards = append(forum.Boards, board)
	}
	forums = append(forums, forum)
	h.renderLayout(w, "index", map[string]interface{}{
		"posts":            "",
		"hasNotifications": false,
		"boards":           boards,
		"motd":             template.HTML(h.motd),
		"logged":           user,
		"title":            settings.Name,
		"forums":           forums,
	}, user)
}
