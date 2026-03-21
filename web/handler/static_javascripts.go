package handler

import (
	"net/http"
	jsEmbed "vpub/assets/js"
)

func (h *Handler) showJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	file, err := jsEmbed.Scripts.ReadFile(r.PathValue("filename"))
	if err != nil {
		serverError(w, err)
	}
	_, err = w.Write(file)
	if err != nil {
		serverError(w, err)
	}
}
