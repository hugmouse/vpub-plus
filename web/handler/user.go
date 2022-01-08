package handler

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"vpub/web/handler/form"
)

func (h *Handler) showUserPostsView(w http.ResponseWriter, r *http.Request) {
	logged, _ := h.session.Get(r)

	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	user, err := h.storage.UserByName(mux.Vars(r)["userId"])
	if err != nil {
		notFound(w)
		return
	}

	//posts, showMore, err := h.storage.PostsByUsernameWithReplyCount(user.Name, h.perPage, page)
	if err != nil {
		serverError(w, err)
		return
	}

	h.renderLayout(w, "user_posts", map[string]interface{}{
		"user":     user,
		"posts":    "",
		"page":     page,
		"showMore": "",
		"nextPage": page + 1,
	}, logged)
}

func (h *Handler) showAccountView(w http.ResponseWriter, r *http.Request, user string) {
	u, err := h.storage.UserByName(user)
	if err != nil {
		forbidden(w)
		return
	}

	h.renderLayout(w, "account", map[string]interface{}{
		"about":          u.About,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) saveAbout(w http.ResponseWriter, r *http.Request, user string) {
	aboutForm := form.NewAboutForm(r)
	if err := h.storage.UpdateAbout(user, aboutForm.About); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/~"+user, http.StatusFound)
}
