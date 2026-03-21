//go:build !debug
// +build !debug

package handler

import "net/http"

func registerDebugHandlers(mux *http.ServeMux) {
	// No-op in non-debug builds
}
