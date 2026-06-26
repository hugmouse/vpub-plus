package handler

import (
	"context"
	"errors"
	"net/http"
	"vpub/storage"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

// setupCompleted reports whether onboarding has already finished. When it has,
// it writes a redirect to "/" and returns true
//
// Essentially replaces the old hardcoded credentials for admin user with a
// nicer interface
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

// setup handles the onboarding form submission
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

	h.setupMu.Lock()
	adminID, settings, err := h.storage.CompleteSetup(storage.SetupRequest{
		AdminName:     setupForm.AdminName,
		AdminPassword: setupForm.AdminPassword,
		ForumName:     setupForm.ForumName,
		URL:           setupForm.URL,
		Lang:          setupForm.Lang,
	})
	h.setupMu.Unlock()

	if err != nil {
		if errors.Is(err, storage.ErrSetupCompleted{}) {
			// Another request finished setup first.
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		v.Set("errorMessage", "Unable to complete setup: "+err.Error())
		v.Render()
		return
	}

	// Hopefully we are safe here
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
