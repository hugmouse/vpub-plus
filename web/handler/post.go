package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"pboard/model"
	"pboard/web/handler/form"
)

func (h *Handler) showNewPostView(w http.ResponseWriter, r *http.Request, user string) {
	form := form.PostForm{}
	form.Topics = h.topics
	h.renderLayout(w, "create_post", map[string]interface{}{
		"form":           form,
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
