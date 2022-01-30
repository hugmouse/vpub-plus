package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"net/url"
	"strconv"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) ParseIntQS(qs *url.URL, name string) (int64, error) {
	if v, ok := qs.Query()[name]; ok && len(v) == 1 {
		return strconv.ParseInt(v[0], 10, 64)
	}
	return 0, errors.New("qs value not found")
}

func (h *Handler) savePost(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.GetUser(r)
	postForm := form.NewPostForm(r)
	post := model.Post{
		User:    user,
		Subject: postForm.Subject,
		Content: postForm.Content,
		TopicId: postForm.TopicId,
	}
	if err := post.Validate(); err != nil {
		serverError(w, err)
		return
	}
	id, err := h.storage.CreatePost(post)
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/topics/%d#%d", postForm.TopicId, id), http.StatusFound)
}

func (h *Handler) showEditPostView(w http.ResponseWriter, r *http.Request) {
	post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
	if err != nil {
		serverError(w, err)
		return
	}
	postForm := form.PostForm{
		Subject: post.Subject,
		Content: post.Content,
		TopicId: post.TopicId,
	}
	h.renderLayout(w, r, "edit_post", map[string]interface{}{
		"form":           postForm,
		"post":           post,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) showEditTopicView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.GetUser(r)
	if !user.IsAdmin {
		notFound(w)
		return
	}
	topic, err := h.storage.TopicById(RouteInt64Param(r, "topicId"))
	if err != nil {
		notFound(w)
		return
	}
	post, err := h.storage.PostById(topic.Post.Id)
	if err != nil {
		notFound(w)
		return
	}
	boards, err := h.storage.Boards()
	topicForm := form.TopicForm{
		Id:       topic.Id,
		BoardId:  topic.BoardId,
		IsSticky: topic.IsSticky,
		IsLocked: topic.IsLocked,
		Boards:   boards,
		PostForm: form.PostForm{
			Subject: post.Subject,
			Content: post.Content,
			TopicId: post.TopicId,
		},
	}
	h.renderLayout(w, r, "edit_topic", map[string]interface{}{
		"form":           topicForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) updateTopic(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.GetUser(r)
	if !user.IsAdmin {
		notFound(w)
		return
	}
	id := RouteInt64Param(r, "topicId")
	t, err := h.storage.TopicById(id)
	if err != nil {
		notFound(w)
		return
	}
	topicForm := form.NewTopicForm(r)
	boardId := topicForm.NewBoardId
	if boardId == 0 {
		boardId = topicForm.BoardId
	}
	topic := model.Topic{
		Id:       id,
		BoardId:  boardId,
		IsSticky: topicForm.IsSticky,
		IsLocked: topicForm.IsLocked,
		Post: model.Post{
			Id:      t.Post.Id,
			Subject: topicForm.PostForm.Subject,
			Content: topicForm.PostForm.Content,
		},
	}
	if err := h.storage.UpdateTopic(topic); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/topics/%d", topic.Id), http.StatusFound)
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.GetUser(r)
	id := RouteInt64Param(r, "postId")
	post, err := h.storage.PostById(id)
	if err != nil {
		serverError(w, err)
		return
	}
	postForm := form.NewPostForm(r)
	post.Subject = postForm.Subject
	post.Content = postForm.Content
	post.User = user
	if err := post.Validate(); err != nil {
		serverError(w, err)
		return
	}
	if err := h.storage.UpdatePost(post); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/topics/%d", post.TopicId), http.StatusFound)
}

func (h *Handler) handleRemovePost(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.GetUser(r)
	switch r.Method {
	case "GET":
		id := RouteInt64Param(r, "postId")
		post, err := h.storage.PostById(id)
		if err != nil {
			serverError(w, err)
			return
		}
		h.renderLayout(w, r, "confirm_remove_post", map[string]interface{}{
			"post":           post,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	case "POST":
		post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
		if err != nil {
			serverError(w, err)
			return
		}
		post.User = user
		err = h.storage.DeletePost(post)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", post.TopicId), http.StatusFound)
	}
}
