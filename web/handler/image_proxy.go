package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TODO: probably rename it since I plan to use it only for images
// and not universal proxy
func (h *Handler) imageProxyHandler(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		http.Error(w, "Invalid URL parameter", http.StatusBadRequest)
		return
	}
	urlStr := parsedURL.String()

	h.cacheMutex.RLock()
	val, ok := h.cachedImages[urlStr]
	isStale := ok && time.Since(val.lastUpdate) > 1*time.Minute
	h.cacheMutex.RUnlock()

	log.Println(h.cachedImages)

	// Cache hit :)
	if ok && !isStale {
		if val.isImage {
			h.serveFromCache(w, val)
			return
		}
		http.Error(w, "Not an image", http.StatusBadGateway)
		return
	}

	// No cache hit :(
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, urlStr, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", urlStr, err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// TODO: move version to global const
	req.Header.Set("User-Agent", "vpub-plus/1.13")

	resp, err := h.httpClient.Do(req)

	// Try to serve from cache upon error
	if err != nil {
		if isStale {
			h.serveFromCache(w, val)
			return
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Received non-cacheable response for %s (Status: %d, Type: %s)",
			urlStr, resp.StatusCode, resp.Header.Get("Content-Type"))

		// Try to serve from cache upon error again
		if isStale {
			log.Printf("Non-cacheable response for %s, serving stale data instead.", urlStr)
			h.serveFromCache(w, val)
			return
		}

		http.Error(w, fmt.Sprintf("Upstream server returned status %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body for %s: %v", urlStr, err)
		if isStale {
			log.Printf("Read body failed for %s, serving stale data: %v", urlStr, err)
			h.serveFromCache(w, val)
			return
		}
		http.Error(w, "Failed to read response body", http.StatusBadGateway)
		return
	}

	isImage := strings.HasPrefix(resp.Header.Get("Content-Type"), "image/")

	var newValue CachedImage

	if isImage {
		newValue = CachedImage{
			lastUpdate: time.Now(),
			isImage:    isImage,
			value:      imageBytes,
		}
	} else {
		newValue = CachedImage{
			isImage: isImage,
		}
	}

	h.cacheMutex.Lock()
	h.cachedImages[urlStr] = newValue
	h.cacheMutex.Unlock()

	log.Printf("Fetched and cached: %s", urlStr)

	h.serveFromCache(w, newValue)
}

func (h *Handler) serveFromCache(w http.ResponseWriter, val CachedImage) {
	expires := val.lastUpdate.Add(1 * time.Minute)
	if time.Now().Before(expires) {
		maxAge := time.Until(expires).Seconds()
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f", maxAge))
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(val.value.([]byte)); err != nil {
		log.Printf("Error writing response body from cache: %v", err)
	}
}
