package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateTopic(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)
	id := RouteInt64Param(r, "topicId")

	topicForm := form.NewTopicForm(r)

	boards, err := h.storage.Boards()
	if err != nil {
		notFound(w)
		return
	}

	topicForm.Boards = boards

	topic, err := h.storage.TopicByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}

	currentBoard, err := h.storage.BoardByID(topic.BoardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}

	if !canAccessForum(currentBoard.Forum, user) {
		forbidden(w)
		return
	}

	board, err := h.storage.BoardByID(topicForm.BoardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}

	if !canAccessForum(board.Forum, user) {
		forbidden(w)
		return
	}

	v := NewView(w, r, "create_topic")
	v.Set("form", topicForm)
	v.Set("board", board)

	boardId := topicForm.NewBoardID
	if boardId == 0 {
		boardId = topicForm.BoardID
	}

	if boardId != topic.BoardID {
		destBoard, err := h.storage.BoardByID(boardId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				notFound(w)
				return
			}
			serverError(w, err)
			return
		}
		if !canAccessForum(destBoard.Forum, user) {
			forbidden(w)
			return
		}
	}

	topicModificationRequest := model.TopicRequest{
		BoardID:  boardId,
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

	err = h.storage.UpdateTopic(id, user.ID, topicModificationRequest)
	if err != nil {
		v.Set("errorMessage", "Unable to updated topic: "+err.Error())
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), http.StatusFound)
}
