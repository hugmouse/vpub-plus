package web

import (
	"fmt"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"pboard/config"
	"pboard/storage"
	"pboard/web/handler"
	"pboard/web/session"
)

func Serve(cfg *config.Config, data *storage.Storage) error {
	var err error
	sess := session.New(cfg.SessionKey, data)
	s, err := handler.New(cfg.Host, cfg.Env, cfg.CSRFKey, data, sess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting HTTP server on %s\n", cfg.Host)
	CSRF := csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(0))
	return http.ListenAndServe(cfg.Host, CSRF(s))
}
