package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showCreateTopicView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}
	boards, err := h.storage.Boards()
	topicForm := form.TopicForm{
		BoardId: board.Id,
		Boards:  boards,
	}
	h.renderLayout(w, r, "create_topic", map[string]interface{}{
		"form":           topicForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
