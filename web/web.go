package web

import (
	"net/http"

	"github.com/holonet/core/logger"
)

// StartServer initializes and starts the HTTP server on the specified address
func StartServer(addr string) {
	logger.Info("Starting web server on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal("Failed to start web server: %v", err)
	}
}
