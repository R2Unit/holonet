package web

import (
	"log"
	"net/http"
)

func StartServer(addr string) {
	http.HandleFunc("/ws", HandleWebSocket)
	log.Printf("WebSocket starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("WebSocket error: %v", err)
	}
}
