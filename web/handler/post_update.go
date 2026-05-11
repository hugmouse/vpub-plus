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

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

	id := RouteInt64Param(r, "postId")

	postForm := form.NewPostForm(r)

	post, err := h.storage.PostByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}

	if postForm.TopicID != 0 && postForm.TopicID != post.TopicID {
		forbidden(w)
		return
	}

	topic, err := h.storage.TopicByID(post.TopicID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}

	board, err := h.storage.BoardByID(topic.BoardID)
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

	v := NewView(w, r, "edit_post")
	v.Set("form", postForm)
	v.Set("board", board)
	v.Set("topic", topic)

	postRequest := model.PostRequest{
		Subject: postForm.Subject,
		Content: postForm.Content,
	}

	if err := validator.ValidatePostRequest(postRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if err := h.storage.UpdatePost(id, user.ID, postRequest); err != nil {
		v.Set("errorMessage", "Unable to create post: "+err.Error())
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d#%d", post.TopicID, id), http.StatusFound)
}
