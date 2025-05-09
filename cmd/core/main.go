package main

import (
	"log"
	_ "net/http"

	"github.com/holonet/core/cache"
	"github.com/holonet/core/database"
	_ "github.com/holonet/core/database/tables"
	"github.com/holonet/core/web"
)

func main() {
	dbHandler, err := database.NewDBHandler()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	log.Printf("Registered %d table(s) for migration.", database.RegisteredTableCount())
	if err := dbHandler.MigrateTables(); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	go dbHandler.StartHeartbeat()

	cacheClient, err := cache.NewCacheClient()
	if err != nil {
		log.Fatalf("Cache initialization error: %v", err)
	}
	log.Println("Cache client connected successfully.")

	go cache.StartHeartbeat(cacheClient)

	go web.StartServer(":3000")

	select {}
}
