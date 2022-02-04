package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/storage"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) checkLogin(w http.ResponseWriter, r *http.Request) {
	loginForm := form.NewLoginForm(r)
	user, err := h.storage.VerifyUser(model.User{
		Name:     loginForm.Username,
		Password: loginForm.Password,
	})
	if err != nil {
		v := NewView(w, r, "login")
		v.Set("form", loginForm)
		switch err.(type) {
		case storage.ErrUserNotFound:
			v.Set("errorMessage", fmt.Sprintf("User %s not found", loginForm.Username))
		case storage.ErrWrongPassword:
			v.Set("errorMessage", "Wrong password")
		}
		v.Render()
		return
	}
	session := request.GetSessionContextKey(r)
	session.SetUserId(user.Id)
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
