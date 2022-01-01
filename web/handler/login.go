package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showLoginView(w http.ResponseWriter, r *http.Request) {
	h.renderLayout(w, "login", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}, "")
}

func (h *Handler) checkLogin(w http.ResponseWriter, r *http.Request) {
	loginForm := form.NewLoginForm(r)
	user, err := h.storage.VerifyUser(model.User{
		Name:     loginForm.Username,
		Password: loginForm.Password,
	})
	if err != nil {
		forbidden(w)
		return
	}
	if err := h.session.Save(r, w, user.Name); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if err := h.session.Delete(w, r); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
