package web

import (
	"log"
	"net/http"
)

// StartServer initializes and starts the HTTP server on the specified address
func StartServer(addr string) {
	log.Printf("Starting web server on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
