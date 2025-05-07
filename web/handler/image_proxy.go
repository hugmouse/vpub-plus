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

const (
	imageSizeLimit = 1024 * 512 // 524KB limit
	cacheTime      = 10 * time.Minute
)

// imageProxyHandler GETs the image from remote host and serves it from a string map
//
// It does not validate the image is in fact a JPEG, PNG or GIF.
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
	isStale := ok && time.Since(val.lastUpdate) > cacheTime
	h.cacheMutex.RUnlock()

	// Cache hit :)
	if ok && !isStale {
		if val.isImage {
			h.serveFromCache(w, val)
			return
		}
		http.Error(w, "Error serving from cache", http.StatusInternalServerError)

		h.cacheMutex.Lock()
		delete(h.cachedImages, urlStr)
		h.cacheMutex.Unlock()

		return
	}

	// No cache hit :(
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", urlStr, err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "vpub-plus/1.13")

	resp, err := h.httpClient.Do(req)

	// Try to serve from cache upon error
	if err != nil {
		if isStale {
			h.serveFromCache(w, val)
			return
		}
		log.Printf("Error getting an image from remote resource %s: %v", urlStr, err)
		http.Error(w, "Error getting an image from remote resource", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to serve from cache upon error again
		if isStale {
			log.Printf("Non-cacheable response for %s, serving stale data instead.", urlStr)
			h.serveFromCache(w, val)
			return
		}

		http.Error(w, fmt.Sprintf("Upstream server returned status %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	// In general, we can use STD's image library to validate JPEG, PNG or GIF,
	// but there's way more image formats that browser supports.
	//
	// So either I have to add more dependencies or rely on remote server not being
	// manipulated such way that they serve __something else__ as an image.
	isImage := strings.HasPrefix(resp.Header.Get("Content-Type"), "image/")

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

	if len(imageBytes) > imageSizeLimit {
		http.Error(w, "Image size too big", http.StatusRequestEntityTooLarge)
		return
	}

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
	expires := val.lastUpdate.Add(cacheTime)
	if time.Now().Before(expires) {
		maxAge := time.Until(expires).Seconds()
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f", maxAge))
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(val.value.([]byte)); err != nil {
		log.Printf("Error writing response body from cache: %v", err)
	}
}
