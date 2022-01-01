package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) saveReplyReply(w http.ResponseWriter, r *http.Request, user string) {
	parent, err := h.storage.ReplyById(RouteInt64Param(r, "replyId"))
	if err != nil {
		serverError(w, err)
		return
	}
	replyForm := form.NewReplyForm(r)
	reply := model.Reply{
		User:     user,
		Content:  replyForm.Content,
		PostId:   parent.PostId,
		ParentId: &parent.Id,
	}
	if _, err := h.storage.CreateReply(reply); err != nil {
		serverError(w, err)
		return
	}
	if parent.User != reply.User {
		if err := h.storage.DeleteNotificationByReplyId(*reply.ParentId); err != nil {
			serverError(w, err)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/replies/%d", *reply.ParentId), http.StatusFound)
}

func (h *Handler) savePostReply(w http.ResponseWriter, r *http.Request, user string) {
	replyForm := form.NewReplyForm(r)
	reply := model.Reply{
		User:    user,
		Content: replyForm.Content,
		PostId:  RouteInt64Param(r, "postId"),
	}
	if err := reply.Validate(); err != nil {
		serverError(w, err)
		return
	}
	if _, err := h.storage.CreateReply(reply); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts/%d", reply.PostId), http.StatusFound)
}

func (h *Handler) showReplyView(w http.ResponseWriter, r *http.Request, user string) {
	reply, err := h.storage.ReplyById(RouteInt64Param(r, "replyId"))
	if err != nil {
		serverError(w, err)
		return
	}
	reply.Thread, err = h.storage.RepliesByParentId(reply.Id)
	if err != nil {
		serverError(w, err)
		return
	}

	post, err := h.storage.PostById(reply.PostId)
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "reply", map[string]interface{}{
		"post":           post,
		"reply":          reply,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showEditReplyView(w http.ResponseWriter, r *http.Request, user string) {
	reply, err := h.storage.ReplyById(RouteInt64Param(r, "replyId"))
	if err != nil {
		serverError(w, err)
		return
	}

	replyForm := form.ReplyForm{
		Content: reply.Content,
	}

	h.renderLayout(w, "edit_reply", map[string]interface{}{
		"reply":          reply,
		"form":           replyForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) updateReply(w http.ResponseWriter, r *http.Request, user string) {
	replyForm := form.NewReplyForm(r)
	reply := model.Reply{
		Id:      RouteInt64Param(r, "replyId"),
		Content: replyForm.Content,
	}
	if err := reply.Validate(); err != nil {
		serverError(w, err)
		return
	}
	if err := h.storage.UpdateReply(reply); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/replies/%d", reply.Id), http.StatusFound)
}

func (h *Handler) handleRemoveReply(w http.ResponseWriter, r *http.Request, user string) {
	reply, err := h.storage.ReplyById(RouteInt64Param(r, "replyId"))
	if err != nil {
		serverError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.renderLayout(w, "confirm_remove_reply", map[string]interface{}{
			"reply":          reply,
			csrf.TemplateTag: csrf.TemplateField(r),
		}, user)
	case http.MethodPost:
		err = h.storage.DeleteReply(reply.Id)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%d", reply.PostId), http.StatusFound)
	}
}
