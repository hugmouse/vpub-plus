//go:build debug
// +build debug

package handler

import (
	_ "expvar"
	"log"
	"net/http"
	"net/http/pprof"
)

func registerDebugHandlers(mux *http.ServeMux) {
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

	debugLog := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("[WARNING] Debug endpoint accessed:", r.URL.Path)
			next(w, r)
		}
	}

	mux.HandleFunc("GET /debug/pprof/", debugLog(pprof.Index))
	mux.HandleFunc("GET /debug/pprof/cmdline", debugLog(pprof.Cmdline))
	mux.HandleFunc("GET /debug/pprof/profile", debugLog(pprof.Profile))
	mux.HandleFunc("GET /debug/pprof/symbol", debugLog(pprof.Symbol))
	mux.HandleFunc("GET /debug/pprof/trace", debugLog(pprof.Trace))
	mux.HandleFunc("GET /debug/pprof/goroutine", debugLog(pprof.Handler("goroutine").ServeHTTP))
	mux.HandleFunc("GET /debug/pprof/heap", debugLog(pprof.Handler("heap").ServeHTTP))
	mux.HandleFunc("GET /debug/pprof/threadcreate", debugLog(pprof.Handler("threadcreate").ServeHTTP))
	mux.HandleFunc("GET /debug/pprof/block", debugLog(pprof.Handler("block").ServeHTTP))
	mux.HandleFunc("GET /debug/vars", debugLog(http.DefaultServeMux.ServeHTTP))
}
