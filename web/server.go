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
	s, err := handler.New(cfg, data, sess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting HTTP server on %s:%s\n", cfg.Host, cfg.Port)
	CSRF := csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(0))
	return http.ListenAndServe(cfg.Host+":"+cfg.Port, CSRF(s))
}
