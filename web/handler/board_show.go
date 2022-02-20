package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) showBoardView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}

	topics, hasMore, err := h.storage.TopicsByBoardId(board.Id, page)
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "board")
	v.Set("navigation", navigation{
		Forum: board.Forum,
		Board: board,
	})
	v.Set("board", board)
	v.Set("topics", topics)
	v.Set("pagination", pagination{
		HasMore: hasMore,
		Page:    page,
	})
	v.Render()
}
