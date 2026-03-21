package handler

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"

	"vpub/web/handler/request"
)

func newCSRFMiddleware(secure bool) func(http.Handler) http.Handler {
	if !secure {
		log.Println("[Warning] CSRF: Secure cookie is disabled. Set CSRF_SECURE=true environment variable to enable it.")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ""
			cookie, err := r.Cookie("csrf_token")
			if err == nil {
				token = cookie.Value
			} else {
				b := make([]byte, 32)
				if _, err := rand.Read(b); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				token = base64.URLEncoding.EncodeToString(b)
				http.SetCookie(w, &http.Cookie{
					Name:     "csrf_token",
					Value:    token,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   86400,
					SameSite: http.SameSiteLaxMode,
					Secure:   secure,
				})
			}

			ctx := context.WithValue(r.Context(), request.CSRFTokenKey, token)
			r = r.WithContext(ctx)

			switch r.Method {
			case "POST", "PUT", "PATCH", "DELETE":
				formToken := r.FormValue("csrf_token")
				if subtle.ConstantTimeCompare([]byte(token), []byte(formToken)) != 1 {
					http.Error(w, "Forbidden - CSRF token invalid", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
