package handler

import (
	"fmt"
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
		v.Set("errorMessage", "Unable to create user")
		v.Render()
		return
	}

	session := request.GetSessionContextKey(r)
	session.SetUserId(id)
	session.FlashInfo(fmt.Sprintf("Welcome, %s!", userCreationRequest.Name))
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}
