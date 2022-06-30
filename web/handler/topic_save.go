package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) saveTopic(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

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

	topicCreationRequest := model.TopicRequest{
		BoardId:  boardId,
		IsSticky: topicForm.IsSticky,
		IsLocked: topicForm.IsLocked,
		Subject:  topicForm.PostForm.Subject,
		Content:  topicForm.PostForm.Content,
	}

	if err := validator.ValidateTopicRequest(topicCreationRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	id, err := h.storage.CreateTopic(user.Id, topicCreationRequest)
	if err != nil {
		v.Set("errorMessage", "Unable to create topic: "+err.Error())
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), http.StatusFound)
}
