package handler

import (
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

	topic, err := h.storage.TopicById(postForm.TopicId)
	if err != nil {
		notFound(w)
		return
	}

	board, err := h.storage.BoardById(topic.BoardId)

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

	if err := h.storage.UpdatePost(id, user.Id, postRequest); err != nil {
		v.Set("errorMessage", "Unable to create post")
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d#%d", postForm.TopicId, id), http.StatusFound)
}
