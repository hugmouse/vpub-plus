package handler

import (
	"errors"
	"net/http"
	"pboard/model"
	"pboard/web/handler/form"
)

func (h *Handler) showRegisterView(w http.ResponseWriter, r *http.Request) {
	h.renderLayout(w, "register", nil, "")
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	userForm := form.NewUserForm(r)
	user := model.User{
		Name:     userForm.Username,
		Password: userForm.Password,
	}
	err := user.Validate()
	if err != nil {
		serverError(w, err)
		return
	}
	if h.storage.UserExists(user.Name) {
		serverError(w, errors.New("username already exists"))
		return
	}
	//if ok := key.Unlock(r.FormValue("key")); !ok {
	//	forbidden(w)
	//	return
	//}
	err = h.storage.CreateUser(user)
	if err != nil {
		serverError(w, err)
		return
	}
	if err := h.session.Save(r, w, r.FormValue("name")); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
