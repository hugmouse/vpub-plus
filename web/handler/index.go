package handler

import (
	"net/http"
	"vpub/model"
)

func forumFromBoards(boards []model.Board) []model.Forum {
	var forums []model.Forum
	var forum model.Forum
	for i, board := range boards {
		if i == 0 {
			forum.Name = board.Forum.Name
			forum.Id = board.Forum.Id
		} else if board.Forum.Id != forum.Id {
			forums = append(forums, forum)
			forum = model.Forum{Name: board.Forum.Name, Id: board.Forum.Id}
		}
		forum.Boards = append(forum.Boards, board)
	}
	if len(forum.Boards) > 0 {
		forums = append(forums, forum)
	}
	return forums
}

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	forums := forumFromBoards(boards)
	h.renderLayout(w, "index", map[string]interface{}{
		"forums": forums,
	}, user)
}
