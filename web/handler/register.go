package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/storage"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) showRegisterView(w http.ResponseWriter, r *http.Request) {
	v := NewView(w, r, "register")
	v.Render()
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	userForm := form.NewUserForm(r)

	v := NewView(w, r, "register")
	v.Set("form", userForm)

	if err := userForm.Validate(); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	user := *userForm.Merge(&model.User{})
	id, err := h.storage.CreateUser(user, userForm.Key)
	if err != nil {
		switch err.(type) {
		case storage.ErrUserExists:
			v.Set("errorMessage", "This username is already taken")
		default:
			v.Set("errorMessage", "An unexpected error happened")
		}
		v.Render()
		return
	}

	session := request.GetSessionContextKey(r)
	session.SetUserId(id)
	session.FlashInfo(fmt.Sprintf("Welcome, %s!", user.Name))
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
