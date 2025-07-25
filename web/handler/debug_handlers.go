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
	log.Println("[handler] - /debug/pprof/vars")
	log.Println("[handler] - To use them, execute: go tool pprof <URL>/debug/pprof/")
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/vars", http.DefaultServeMux)
}
