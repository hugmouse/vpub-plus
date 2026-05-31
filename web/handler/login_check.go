package handler

import (
	"context"
	"errors"
	"net/http"
	"vpub/model"
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
		v.Set("errorMessage", "Invalid username or password")
		v.Render()
		return
	}

	newSession, err := h.session.NewSession(w, r, user.ID)
	if err != nil {
		serverError(w, errors.New("can't create a new session: "+err.Error()))
		return
	}

	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}

	ctx := context.WithValue(r.Context(), request.SessionKey, newSession)
	ctx = context.WithValue(ctx, request.UserKey, user)
	ctx = context.WithValue(ctx, request.SettingsKey, settings)

	http.Redirect(w, r.WithContext(ctx), "/", http.StatusFound)
}
