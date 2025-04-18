package api

import (
	"log"
	"net/http"

	"github.com/holonet/core/database"
)

var dbHandler *database.DBHandler

func SetDBHandler(handler *database.DBHandler) {
	dbHandler = handler
}

func StartServer(addr string) {
	log.Printf("Starting API service on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting API service: %v", err)
	}
}
