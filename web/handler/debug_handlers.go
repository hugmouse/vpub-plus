//go:build debug
// +build debug

package handler

import (
	_ "expvar"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

func registerDebugHandlers(router *mux.Router) {
	log.Println("[handler] Debug routes added:")
	log.Println("[handler] - /debug/pprof/")
	log.Println("[handler] - /debug/pprof/cmdline")
	log.Println("[handler] - /debug/pprof/profile")
	log.Println("[handler] - /debug/pprof/symbol")
	log.Println("[handler] - /debug/pprof/trace")
	log.Println("[handler] - /debug/pprof/goroutine")
	log.Println("[handler] - /debug/pprof/heap")
	log.Println("[handler] - /debug/pprof/threadcreate")
	log.Println("[handler] - /debug/pprof/block")
	log.Println("[handler] - /debug/vars")
	log.Println("[handler] - To use them, execute: go tool pprof <URL>/debug/pprof/")

	debugRouter := router.PathPrefix("/debug").Subrouter()
	debugRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("[WARNING] Debug endpoint accessed:", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	debugRouter.HandleFunc("/pprof/", pprof.Index)
	debugRouter.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	debugRouter.HandleFunc("/pprof/profile", pprof.Profile)
	debugRouter.HandleFunc("/pprof/symbol", pprof.Symbol)
	debugRouter.HandleFunc("/pprof/trace", pprof.Trace)
	debugRouter.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	debugRouter.Handle("/pprof/heap", pprof.Handler("heap"))
	debugRouter.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	debugRouter.Handle("/pprof/block", pprof.Handler("block"))
	debugRouter.Handle("/vars", http.DefaultServeMux)
}
