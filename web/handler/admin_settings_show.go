package handler

import (
	"github.com/gorilla/csrf"
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
	}
	h.renderLayout(w, r, "admin_settings_edit", map[string]interface{}{
		"form":           settingsForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
