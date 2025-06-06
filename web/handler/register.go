package handler

import (
	"context"
	"errors"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	userForm := form.NewUserForm(r)

	v := NewView(w, r, "register")
	v.Set("form", userForm)

	if err := userForm.Validate(); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	userCreationRequest := model.UserCreationRequest{
		Name:     userForm.Username,
		Password: userForm.Password,
	}

	if err := validator.ValidateUserCreation(h.storage, userForm.Key, userCreationRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	id, err := h.storage.CreateUser(userForm.Key, userCreationRequest)
	if err != nil {
		v.Set("errorMessage", "Unable to create user: "+err.Error())
		v.Render()
		return
	}

	newSession, err := h.session.NewSession(w, r, id)
	if err != nil {
		serverError(w, errors.New("can't create a new session: "+err.Error()))
		return
	}

	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, request.SessionKey, newSession)
	ctx = context.WithValue(ctx, request.SettingsKey, settings)

	// TODO: fix flash
	// session.FlashInfo(fmt.Sprintf("Welcome, %s!", userCreationRequest.Name))

	http.Redirect(w, r, "/", http.StatusFound)
}
