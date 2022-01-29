package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/model"
	"vpub/storage"
	"vpub/web/handler/form"
)

func (h *Handler) showLoginView(w http.ResponseWriter, r *http.Request) {
	h.renderLayoutFlash(w, r, "login", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}, model.User{})
}

func (h *Handler) checkLogin(w http.ResponseWriter, r *http.Request) {
	session, err := h.session.GetSession(r)
	if err != nil {
		serverError(w, err)
		return
	}
	loginForm := form.NewLoginForm(r)
	user, err := h.storage.VerifyUser(model.User{
		Name:     loginForm.Username,
		Password: loginForm.Password,
	})
	if err != nil {
		switch err.(type) {
		case storage.ErrUserNotFound:
			session.AddFlash(fmt.Sprintf("User %s not found", loginForm.Username), "errors")
		case storage.ErrWrongPassword:
			session.AddFlash("Wrong password", "errors")
		}
		err := session.Save(r, w)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	session.Values["id"] = user.Id
	if err := session.Save(r, w); err != nil {
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
