package main

import (
	_ "net/http"

	"github.com/holonet/core/cache"
	"github.com/holonet/core/database"
	_ "github.com/holonet/core/database/tables"
	"github.com/holonet/core/logger"
	"github.com/holonet/core/web"
)

func main() {
	// Initialize logger with command line flags
	logger.Init()

	dbHandler, err := database.NewDBHandler()
	if err != nil {
		logger.Fatal("Failed to initialize DB: %v", err)
	}

	logger.Info("Registered %d table(s) for migration.", database.RegisteredTableCount())
	if err := dbHandler.MigrateTables(); err != nil {
		logger.Fatal("Migration error: %v", err)
	}
	logger.Info("Database migrations completed successfully.")

	go dbHandler.StartHeartbeat()

	cacheClient, err := cache.NewCacheClient()
	if err != nil {
		logger.Fatal("Cache initialization error: %v", err)
	}
	logger.Info("Cache client connected successfully.")

	go cache.StartHeartbeat(cacheClient)

	go web.StartServer(":3000")

	logger.Debug("Main goroutine waiting indefinitely")
	select {}
}
