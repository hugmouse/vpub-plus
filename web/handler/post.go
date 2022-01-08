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

func (h *Handler) showNewPostView(w http.ResponseWriter, r *http.Request, user string) {
	id, err := h.ParseIntQS(r.URL, "topicId")
	if err != nil {
		notFound(w)
		return
	}
	topic, err := h.storage.TopicById(id)
	if err != nil {
		notFound(w)
		return
	}
	postForm := form.PostForm{}
	postForm.Topics = h.topics
	postForm.Topic = topic.Name
	h.renderLayout(w, "create_post", map[string]interface{}{
		"form":           postForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) savePost(w http.ResponseWriter, r *http.Request, user string) {
	postForm := form.NewPostForm(r, h.topics)
	post := model.Post{
		User:    user,
		Title:   postForm.Title,
		Content: postForm.Content,
		Topic:   postForm.Topic,
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
	http.Redirect(w, r, fmt.Sprintf("/posts/%d", id), http.StatusFound)
}

func (h *Handler) showPostView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
	if err != nil {
		notFound(w)
		return
	}

	replies, err := h.storage.RepliesByPostId(post.Id)
	if err != nil {
		serverError(w, err)
		return
	}

	h.renderLayout(w, "post", map[string]interface{}{
		"post":           post,
		"content":        post.Content,
		"replies":        replies,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showEditPostView(w http.ResponseWriter, r *http.Request, user string) {
	post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
	if err != nil {
		serverError(w, err)
		return
	}
	if user != post.User {
		forbidden(w)
		return
	}
	postForm := form.PostForm{
		Title:   post.Title,
		Content: post.Content,
		Topic:   post.Topic,
		Topics:  h.topics,
	}
	h.renderLayout(w, "edit_post", map[string]interface{}{
		"form":           postForm,
		"post":           post,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request, user string) {
	id := RouteInt64Param(r, "postId")

	post, err := h.storage.PostById(id)
	if err != nil {
		serverError(w, err)
		return
	}

	if user != post.User {
		forbidden(w)
		return
	}

	postForm := form.NewPostForm(r, h.topics)

	post.Title = postForm.Title
	post.Content = postForm.Content
	post.Topic = postForm.Topic
	post.User = user

	if err := post.Validate(); err != nil {
		serverError(w, err)
		return
	}

	if err := h.storage.UpdatePost(post); err != nil {
		serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%d", post.Id), http.StatusFound)
}

func (h *Handler) handleRemovePost(w http.ResponseWriter, r *http.Request, user string) {
	switch r.Method {
	case "GET":
		id := RouteInt64Param(r, "postId")
		post, err := h.storage.PostById(id)
		if err != nil {
			serverError(w, err)
			return
		}
		if user != post.User {
			forbidden(w)
			return
		}
		h.renderLayout(w, "confirm_remove_post", map[string]interface{}{
			"post":           post,
			csrf.TemplateTag: csrf.TemplateField(r),
		}, user)
	case "POST":
		post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
		if err != nil {
			serverError(w, err)
			return
		}
		if user != post.User {
			forbidden(w)
			return
		}
		err = h.storage.DeletePost(post.Id, user)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
