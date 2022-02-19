package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
)

func (h *Handler) updateTopic(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "topicId")

	topicForm := form.NewTopicForm(r)

	boards, err := h.storage.Boards()
	if err != nil {
		notFound(w)
		return
	}

	topicForm.Boards = boards

	board, err := h.storage.BoardById(topicForm.BoardId)
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "create_topic")
	v.Set("form", topicForm)
	v.Set("board", board)

	boardId := topicForm.NewBoardId
	if boardId == 0 {
		boardId = topicForm.BoardId
	}

	topicModificationRequest := model.TopicRequest{
		BoardId:  boardId,
		IsSticky: topicForm.IsSticky,
		IsLocked: topicForm.IsLocked,
		Subject:  topicForm.PostForm.Subject,
		Content:  topicForm.PostForm.Content,
	}

	if err := validator.ValidateTopicRequest(topicModificationRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	err = h.storage.UpdateTopic(id, topicModificationRequest)
	if err != nil {
		v.Set("errorMessage", "Unable to updated topic")
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), http.StatusFound)
}
