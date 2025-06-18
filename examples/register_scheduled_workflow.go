package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/holonet/core/workflow"
	_ "github.com/lib/pq"
)

// This script demonstrates how to register and schedule the housekeeping workflow
func main() {
	// Connect to the database
	db, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a workflow manager
	workflowManager := workflow.NewWorkflowManager(db)

	// Register the housekeeping workflow
	housekeepingWorkflow, err := registerHousekeepingWorkflow(workflowManager)
	if err != nil {
		log.Fatalf("Failed to register housekeeping workflow: %v", err)
	}

	// Schedule the workflow to run daily at midnight
	scheduledTime := time.Now().Add(24 * time.Hour)
	scheduledTime = time.Date(
		scheduledTime.Year(),
		scheduledTime.Month(),
		scheduledTime.Day(),
		0, 0, 0, 0,
		scheduledTime.Location(),
	)

	// Create parameters for the workflow
	parameters := map[string]interface{}{
		"max_age_days": 30,
		"dry_run":      false,
	}
	parametersJSON, err := json.Marshal(parameters)
	if err != nil {
		log.Fatalf("Failed to marshal parameters: %v", err)
	}

	// Schedule the workflow
	execution, err := workflowManager.ScheduleWorkflow(
		housekeepingWorkflow.ID,
		parametersJSON,
		scheduledTime,
	)
	if err != nil {
		log.Fatalf("Failed to schedule workflow: %v", err)
	}

	fmt.Printf("Housekeeping workflow scheduled successfully!\n")
	fmt.Printf("Workflow ID: %d\n", housekeepingWorkflow.ID)
	fmt.Printf("Execution ID: %d\n", execution.ID)
	fmt.Printf("Scheduled time: %s\n", scheduledTime.Format(time.RFC3339))
}

// registerHousekeepingWorkflow registers the housekeeping workflow
func registerHousekeepingWorkflow(manager *workflow.WorkflowManager) (*workflow.Workflow, error) {
	// Define the workflow code
	workflowCode := `
package main

import (
	"fmt"
	"time"
)

// HousekeepingWorkflow is an example workflow that performs basic housekeeping tasks
func main() {
	fmt.Println("Starting housekeeping workflow")
	
	// Simulate housekeeping tasks
	cleanupTempFiles()
	cleanupOldLogs()
	performSystemChecks()
	
	fmt.Println("Housekeeping workflow completed successfully")
}

func cleanupTempFiles() {
	fmt.Println("Cleaning up temporary files...")
	// In a real implementation, this would delete old temporary files
	time.Sleep(500 * time.Millisecond) // Simulate work
	fmt.Println("Temporary files cleaned up")
}

func cleanupOldLogs() {
	fmt.Println("Cleaning up old log files...")
	// In a real implementation, this would archive or delete old log files
	time.Sleep(500 * time.Millisecond) // Simulate work
	fmt.Println("Old log files cleaned up")
}

func performSystemChecks() {
	fmt.Println("Performing system health checks...")
	// In a real implementation, this would check system health metrics
	time.Sleep(500 * time.Millisecond) // Simulate work
	fmt.Println("System health checks completed")
}
`

	// Create the workflow
	wf, err := manager.CreateWorkflow(
		"Housekeeping (Scheduled)",
		"Performs routine housekeeping tasks like cleaning up temporary files and logs (scheduled version)",
		workflowCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// Set the workflow to active
	wf.Status = workflow.StatusActive
	if err := manager.UpdateWorkflow(wf); err != nil {
		return nil, fmt.Errorf("failed to activate workflow: %w", err)
	}

	fmt.Printf("Housekeeping workflow registered with ID: %d\n", wf.ID)
	return wf, nil
}

// connectToDatabase connects to the PostgreSQL database
func connectToDatabase() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=holonet password=insecure dbname=holonet sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
