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

func (h *Handler) showNewThreadView(w http.ResponseWriter, r *http.Request, user string) {
	id, err := h.ParseIntQS(r.URL, "topicId")
	if err != nil {
		notFound(w)
		return
	}
	topic, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}
	postForm := form.ThreadForm{
		Topic: topic,
	}
	h.renderLayout(w, "create_post", map[string]interface{}{
		"form":           postForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) savePost(w http.ResponseWriter, r *http.Request, user string) {
	postForm := form.NewPostForm(r)
	post := model.Post{
		User:    user,
		Title:   postForm.Subject,
		Content: postForm.Content,
		TopicId: postForm.TopicId,
	}
	if err := post.Validate(); err != nil {
		serverError(w, err)
		return
	}
	_, err := h.storage.CreatePost(post)
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/topics/%d", postForm.TopicId), http.StatusFound)
}

//
//func (h *Handler) showPostView(w http.ResponseWriter, r *http.Request) {
//	user, _ := h.session.Get(r)
//
//	post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
//	if err != nil {
//		notFound(w)
//		return
//	}
//
//	replies, err := h.storage.RepliesByPostId(post.Id)
//	if err != nil {
//		serverError(w, err)
//		return
//	}
//
//	h.renderLayout(w, "post", map[string]interface{}{
//		"post":           post,
//		"content":        post.Content,
//		"replies":        replies,
//		csrf.TemplateTag: csrf.TemplateField(r),
//	}, user)
//}

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
		Subject: post.Title,
		Content: post.Content,
		TopicId: post.TopicId,
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
	postForm := form.NewPostForm(r)
	post.Title = postForm.Subject
	post.Content = postForm.Content
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
		err = h.storage.DeletePostById(post.Id)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", post.TopicId), http.StatusFound)
	}
}
