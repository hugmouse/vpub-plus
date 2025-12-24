package handler

import (
	"log"
	"net/http"
	"time"
	"vpub/web/handler/request"
)

func (h *Handler) showAdminImageCache(w http.ResponseWriter, r *http.Request) {
	settings := request.GetSettingsContextKey(r)
	cacheTime := time.Duration(settings.ImageProxyCacheTime) * time.Second

	v := NewView(w, r, "admin_image_cache")
	v.Set("images", h.imageProxy.List(cacheTime))
	v.Render()
}

func (h *Handler) removeAdminImageCache(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	if r.FormValue("action") == "all" {
		h.imageProxy.RemoveAll()
		log.Println("Admin cleared image cache")
	} else {
		urlStr := r.FormValue("url")
		if urlStr != "" {
			h.imageProxy.Remove(urlStr)
			log.Printf("Admin removed image from cache: %s", urlStr)
		}
	}

	http.Redirect(w, r, "/admin/image-proxy", http.StatusFound)
}
