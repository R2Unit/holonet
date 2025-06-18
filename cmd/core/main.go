package main

import (
	"context"
	_ "net/http"

	"github.com/holonet/core/api"
	"github.com/holonet/core/cache"
	"github.com/holonet/core/database"
	_ "github.com/holonet/core/database/tables"
	"github.com/holonet/core/logger"
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

	api.RegisterWorkflowRoutes()

	api.RegisterPolicyRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go workflowExecutor.StartExecutionLoop(ctx)

	go web.StartServer(":3000")

	logger.Debug("Main goroutine waiting indefinitely")
	select {}
}
