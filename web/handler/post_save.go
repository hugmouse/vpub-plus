package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) savePost(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

	postForm := form.NewPostForm(r)

	topic, err := h.storage.TopicByID(postForm.TopicID)
	if err != nil {
		notFound(w)
		return
	}

	board, err := h.storage.BoardByID(topic.BoardID)
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "create_post")
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

	id, err := h.storage.CreatePost(user.ID, postForm.TopicID, postRequest)
	if err != nil {
		v.Set("errorMessage", "Unable to create post: "+err.Error())
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d#%d", postForm.TopicID, id), http.StatusFound)
}
