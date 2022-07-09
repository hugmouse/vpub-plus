package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminSettingsView(w http.ResponseWriter, r *http.Request) {
	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}

	settingsForm := form.SettingsForm{
		Name:    settings.Name,
		Css:     settings.Css,
		Footer:  settings.Footer,
		PerPage: settings.PerPage,
		URL:     settings.URL,
		Lang:    settings.Lang,
	}

	v := NewView(w, r, "admin_settings_edit")
	v.Set("form", settingsForm)
	v.Render()
}
