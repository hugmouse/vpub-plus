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
	h.renderLayout(w, "create_post", map[string]interface{}{
		"form":           form,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) savePost(w http.ResponseWriter, r *http.Request, user string) {
	postForm := form.NewPostForm(r)
	post := model.Post{
		User:    user,
		Title:   postForm.Title,
		Content: postForm.Content,
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
