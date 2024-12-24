package handler

import (
	"errors"
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
		if errors.As(err, &storage.ErrUserExists{}) {
			v.Set("errorMessage", fmt.Sprintf("User %s not found", loginForm.Username))
		} else if errors.As(err, &storage.ErrWrongPassword{}) {
			v.Set("errorMessage", "Wrong password")
		} else {
			v.Set("errorMessage", fmt.Errorf("unknown error occurred: %w", err))
		}
		v.Render()
		return
	}
	session := request.GetSessionContextKey(r)
	session.SetUserId(user.Id)
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
