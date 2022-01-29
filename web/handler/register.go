package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/model"
	"vpub/storage"
	"vpub/web/handler/form"
)

func (h *Handler) showRegisterView(w http.ResponseWriter, r *http.Request) {
	h.renderLayoutFlash(w, r, "register", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}, model.User{})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	session, err := h.session.GetSession(r)
	if err != nil {
		serverError(w, err)
		return
	}
	userForm := form.NewUserForm(r)
	user := model.User{
		Name:     userForm.Username,
		Password: userForm.Password,
	}
	if err := userForm.Validate(); err != nil {
		session.AddFlash(err.Error(), "errors")
		if err := session.Save(r, w); err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	id, err := h.storage.CreateUser(user, userForm.Key)
	if err != nil {
		switch err.(type) {
		case storage.ErrUserExists:
			session.AddFlash("This username is already taken", "errors")
		default:
			serverError(w, err)
		}
		if err := session.Save(r, w); err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	session.Values["id"] = id
	session.AddFlash(fmt.Sprintf("Welcome, %s!", user.Name), "info")
	if err := session.Save(r, w); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
