//go:build !debug
// +build !debug

// Provides a stub for registerDebugHandlers when the "debug" build tag is not present
package handler

import "github.com/gorilla/mux"

func registerDebugHandlers(router *mux.Router) {
	// No-op in non-debug builds
}
