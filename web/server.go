package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/csrf"

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
	log.Printf("Starting HTTP server on localhost:%s\n", cfg.Port)

	// MaxAge(86400) persists the CSRF cookie for 24h (vs session-only with 0).
	// This improves usability: forms survive browser restarts.
	var CSRF func(http.Handler) http.Handler
	if cfg.CSRFSecure {
		CSRF = csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(86400), csrf.Path("/"))
	} else {
		log.Println("[Warning] CSRF: Secure cookie is disabled. Set CSRF_SECURE=true environment variable to enable it.")
		CSRF = csrf.Protect([]byte(cfg.CSRFKey), csrf.MaxAge(86400), csrf.Path("/"), csrf.Secure(false))
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      CSRF(s),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-quit:
	}

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}
