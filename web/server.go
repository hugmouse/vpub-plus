package web

import (
	"log"
	"net/http"
	"vpub/config"
	"vpub/storage"
	"vpub/web/handler"
	"vpub/web/session"

	"github.com/gorilla/csrf"
)

func Serve(cfg *config.Config, data *storage.Storage) error {
	var err error
	sess := session.New(cfg.SessionKey, data)
	s, err := handler.New(data, sess)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting HTTP server on localhost:%s\n", cfg.Port)
	var CSRF func(http.Handler) http.Handler
	if cfg.CSRFSecure {
		CSRF = csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(0), csrf.Path("/"))
	} else {
		log.Println("[Warning] CSRF: Secure cookie is disabled. Set CSRF_SECURE=true environment variable to enable it.")
		CSRF = csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(0), csrf.Path("/"), csrf.Secure(false))
	}
	return http.ListenAndServe(":"+cfg.Port, CSRF(s))
}
