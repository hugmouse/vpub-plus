package handler

import (
	"context"
	"errors"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

// setupCompleted reports whether onboarding has already finished. When it has,
// it writes a redirect to "/" and returns true.
//
// Essentially replaces the old hardcoded credentials for admin user with a
// nicer interface.
func (h *Handler) setupCompleted(w http.ResponseWriter, r *http.Request) bool {
	adminExists, err := h.storage.HasAdmin()
	if err != nil {
		serverError(w, err)
		return true
	}
	if adminExists {
		http.Redirect(w, r, "/", http.StatusFound)
		return true
	}
	return false
}

func (h *Handler) showSetupView(w http.ResponseWriter, r *http.Request) {
	if h.setupCompleted(w, r) {
		return
	}

	v := NewView(w, r, "setup")
	v.Set("form", form.SetupForm{})
	v.Render()
}

func (h *Handler) setup(w http.ResponseWriter, r *http.Request) {
	if h.setupCompleted(w, r) {
		return
	}

	setupForm := form.NewSetupForm(r)

	v := NewView(w, r, "setup")
	v.Set("form", setupForm)

	if err := setupForm.Validate(); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	adminID, err := h.storage.CreateAdminUser(model.UserCreationRequest{
		Name:     setupForm.AdminName,
		Password: setupForm.AdminPassword,
		IsAdmin:  true,
	})
	if err != nil {
		v.Set("errorMessage", "Unable to create admin user: "+err.Error())
		v.Render()
		return
	}

	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}
	settings.Name = setupForm.ForumName
	settings.URL = setupForm.URL
	settings.Lang = setupForm.Lang
	if err := h.storage.UpdateSettings(settings); err != nil {
		serverError(w, err)
		return
	}

	if err := h.seedSampleContent(adminID); err != nil {
		serverError(w, err)
		return
	}

	h.setupComplete.Store(true)

	newSession, err := h.session.NewSession(w, r, adminID)
	if err != nil {
		serverError(w, errors.New("can't create a new session: "+err.Error()))
		return
	}

	ctx := context.WithValue(r.Context(), request.SessionKey, newSession)
	ctx = context.WithValue(ctx, request.SettingsKey, settings)
	http.Redirect(w, r.WithContext(ctx), "/", http.StatusFound)
}

// seedSampleContent creates a welcome post for admin user.
// Called right after setup is complete.
func (h *Handler) seedSampleContent(userID int64) error {
	forumID, err := h.storage.CreateForum(model.ForumRequest{
		Name:     "Getting started",
		Position: 0,
		IsLocked: false,
	})
	if err != nil {
		return err
	}

	boardID, err := h.storage.CreateBoard(model.BoardRequest{
		Name:        "Your first board",
		Description: "This is a sample board. Feel free to edit or remove it from the admin area.",
		Position:    0,
		IsLocked:    false,
		ForumID:     forumID,
	})
	if err != nil {
		return err
	}

	_, err = h.storage.CreateTopic(userID, model.TopicRequest{
		Subject: "Welcome to your new forum",
		Content: `Welcome! You just finished the initial setup.

## What now?

- Head to the [admin area](/admin) to manage forums, boards, users and instance settings.
- Create new boards under **Getting started** or remove this sample content entirely.
- Adjust the forum name, URL and appearance under [instance settings](/admin/settings/edit).

Enjoy your forum!`,
		BoardID:  boardID,
		IsLocked: true,
		IsSticky: true,
	})
	return err
}
