package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"vpub/model"
	"vpub/web/handler/request"
)

var (
	darkAndEmptyPNG = [67]uint8{
		// Offset 0x00000000 to 0x00000042
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x25,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x2C, 0x62, 0x57, 0x02, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18, 0x05, 0x00, 0x01,
		0x03, 0x00, 0x01, 0x45, 0x41, 0x03, 0xB5, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
)

func NewImageProxyHandler() *ImageProxyHandler {
	return &ImageProxyHandler{
		cachedImages: make(map[string]CachedImage),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				TLSHandshakeTimeout:   3 * time.Second,
				IdleConnTimeout:       10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func TestImageProxyHandler_ServeHTTP_WithStdHttp(t *testing.T) {
	handlerForImageProxy := NewImageProxyHandler()
	darkAndEmptyPNGSlice := darkAndEmptyPNG[:]

	mockImageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/dark.png":
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(darkAndEmptyPNGSlice)
		case "/not-image":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(darkAndEmptyPNGSlice[7:67])
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockImageServer.Close()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/proxy", handlerForImageProxy.imageProxyHandler)

	tests := []struct {
		name       string
		urlPath    string
		targetURL  string
		wantStatus int
		wantHeader map[string]string
		wantBody   []byte
	}{
		{
			name:       "Proxifying darkAndEmptyPNG via query parameter",
			urlPath:    "/proxy?url=" + url.QueryEscape(mockImageServer.URL+"/dark.png"),
			targetURL:  mockImageServer.URL + "/dark.png",
			wantStatus: http.StatusOK,
			wantHeader: map[string]string{"Content-Type": "image/png"},
			wantBody:   darkAndEmptyPNGSlice,
		},
		{
			name:       "Non-existent image via query parameter",
			urlPath:    "/proxy?url=" + url.QueryEscape(mockImageServer.URL+"/doesnotexist.jpg"),
			targetURL:  mockImageServer.URL + "/doesnotexist.jpg",
			wantStatus: http.StatusBadGateway,
		},
		{
			name:       "Invalid image url via query parameter",
			urlPath:    "/proxy?url=" + url.QueryEscape(mockImageServer.URL+"/not-image"),
			targetURL:  mockImageServer.URL + "/not-image",
			wantStatus: http.StatusBadGateway,
		},
		{
			name:       "Missing URL parameter",
			urlPath:    "/proxy",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid URL",
			urlPath:    "/proxy?url=%invalid%",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Unsupported URL scheme",
			urlPath:    "/proxy?url=ftp://example.com/image.jpg",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			// Inject key settings for test
			settings := model.Settings{
				ImageProxySizeLimit: 1024 * 1024, // 1MB
				ImageProxyCacheTime: 60,
			}
			ctx := context.WithValue(req.Context(), request.SettingsKey, settings)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			serveMux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("ServeHTTP() status code = %+v, want %+v for URL: %s", rr.Code, tt.wantStatus, tt.urlPath)
			}

			for key, value := range tt.wantHeader {
				if got := rr.Header().Get(key); got != value {
					t.Errorf("ServeHTTP() header[%s] = %+v, want %+v for URL: %s", key, got, value, tt.urlPath)
				}
			}

			if len(tt.wantBody) > 0 {
				if !bytes.Equal(rr.Body.Bytes(), tt.wantBody) {
					t.Errorf("ServeHTTP() body mismatch for URL: %s. Got %d bytes, want %d bytes.", tt.urlPath, len(rr.Body.Bytes()), len(tt.wantBody))
				}
			}
		})
	}
}

func TestImageProxyHandler_Management(t *testing.T) {
	h := NewImageProxyHandler()
	urlStr := "http://example.com/image.png"

	// Add an item manually to cache
	h.cacheMutex.Lock()
	h.cachedImages[urlStr] = CachedImage{
		lastUpdate:  time.Now(),
		value:       []byte("fake-image-data"),
		contentType: "image/png",
	}
	h.cacheMutex.Unlock()

	// Test List
	list := h.List(time.Minute)
	if len(list) != 1 {
		t.Errorf("List() returned %d items, want 1", len(list))
	}
	if list[0].URL != urlStr {
		t.Errorf("List() item URL = %s, want %s", list[0].URL, urlStr)
	}

	// Test Remove
	h.Remove(urlStr)
	list = h.List(time.Minute)
	if len(list) != 0 {
		t.Errorf("List() returned %d items after Remove, want 0", len(list))
	}

	// Add again for RemoveAll
	h.cacheMutex.Lock()
	h.cachedImages[urlStr] = CachedImage{
		lastUpdate:  time.Now(),
		value:       []byte("fake-image-data"),
		contentType: "image/png",
	}
	h.cacheMutex.Unlock()

	// Test RemoveAll
	h.RemoveAll()
	list = h.List(time.Minute)
	if len(list) != 0 {
		t.Errorf("List() returned %d items after RemoveAll, want 0", len(list))
	}
}
