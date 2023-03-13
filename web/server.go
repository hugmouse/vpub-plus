package web

import (
	"fmt"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"vpub/config"
	"vpub/storage"
	"vpub/web/handler"
	"vpub/web/session"
)

func Serve(cfg *config.Config, data *storage.Storage) error {
	var err error
	sess := session.New(cfg.SessionKey, data)
	s, err := handler.New(data, sess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting HTTP server on localhost:%s\n", cfg.Port)
	var CSRF func(http.Handler) http.Handler
	if cfg.CSRFSecure {
		fmt.Println("INFO: CSRF: Secure cookie is disabled")
		CSRF = csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(0), csrf.Path("/"), csrf.Secure(false))
	} else {
		CSRF = csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(0), csrf.Path("/"))
	}
	return http.ListenAndServe(":"+cfg.Port, CSRF(s))
}
