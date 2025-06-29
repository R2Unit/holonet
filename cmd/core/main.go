package main

import (
	"context"
	_ "net/http"

	"github.com/holonet/core/api"
	"github.com/holonet/core/api/users"
	"github.com/holonet/core/cache"
	"github.com/holonet/core/database"
	_ "github.com/holonet/core/database/tables"
	"github.com/holonet/core/logger"
	"github.com/holonet/core/netbox"
	"github.com/holonet/core/web"
	"github.com/holonet/core/workflow"
)

func main() {
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
	// Initialize admin user, group, and token if they don't exist
	if err := database.InitAdminUser(dbHandler.DB); err != nil {
		logger.Fatal("Failed to initialize admin user: %v", err)
	}

	go dbHandler.StartHeartbeat()

	cacheClient, err := cache.NewCacheClient()
	if err != nil {
		logger.Fatal("Cache initialization error: %v", err)
	}
	logger.Info("Cache client connected successfully.")

	go cache.StartHeartbeat(cacheClient)

	if err := workflow.Init(dbHandler.DB); err != nil {
		logger.Fatal("Failed to initialize workflow system: %v", err)
	}
	workflowManager := workflow.NewWorkflowManager(dbHandler.DB)
	workflowExecutor := workflow.NewExecutor(workflowManager)

	api.SetDBHandler(dbHandler)
	api.SetWorkflowManager(workflowManager)
	users.SetDBHandler(dbHandler.DB)

	api.RegisterEndpoints()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go workflowExecutor.StartExecutionLoop(ctx)

	go web.StartServer(":3000")

	// Initialize NetBox client and gatekeeper after all other initializations
	netboxClient, err := netbox.New(nil)
	if err != nil {
		logger.Error("Failed to initialize NetBox client: %v", err)
	} else {
		logger.Info("NetBox client initialized successfully.")

		// Initialize the gatekeeper for API request management
		gatekeeper := netbox.NewGatekeeper(netboxClient)
		logger.Info("NetBox gatekeeper initialized successfully.")

		// Initialize NetBox authentication
		if err := netbox.InitNetboxAuth(netboxClient, gatekeeper, dbHandler.DB); err != nil {
			logger.Error("Failed to initialize NetBox authentication: %v", err)
		} else {
			logger.Info("NetBox authentication initialized successfully.")
		}

		// Start the heartbeat to monitor NetBox availability
		netbox.StartHeartbeat(netboxClient)
		logger.Info("NetBox heartbeat started successfully.")
	}

	logger.Debug("Main goroutine waiting indefinitely")
	select {}
}
