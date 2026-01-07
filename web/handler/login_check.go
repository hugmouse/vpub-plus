package handler

import (
	"context"
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

	newSession, err := h.session.NewSession(w, r, user.Id)
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
