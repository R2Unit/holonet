package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/holonet/core/workflow"
	_ "github.com/lib/pq"
)

// This script demonstrates how to register and directly trigger the housekeeping workflow via API
func main() {
	// Connect to the database
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a workflow manager
	workflowManager := workflow.NewWorkflowManager(db)

	// Register the housekeeping workflow
	housekeepingWorkflow, err := registerDirectHousekeepingWorkflow(workflowManager)
	if err != nil {
		log.Fatalf("Failed to register housekeeping workflow: %v", err)
	}

	// Create parameters for the workflow
	parameters := map[string]interface{}{
		"max_age_days": 30,
		"dry_run":      false,
	}
	parametersJSON, err := json.Marshal(parameters)
	if err != nil {
		log.Fatalf("Failed to marshal parameters: %v", err)
	}

	// Trigger the workflow immediately via API
	err = triggerDirectWorkflowViaAPI(housekeepingWorkflow.ID, parametersJSON)
	if err != nil {
		log.Fatalf("Failed to trigger workflow via API: %v", err)
	}

	fmt.Printf("Housekeeping workflow triggered successfully via API!\n")
	fmt.Printf("Workflow ID: %d\n", housekeepingWorkflow.ID)
}

// registerDirectHousekeepingWorkflow registers the housekeeping workflow for direct triggering
func registerDirectHousekeepingWorkflow(manager *workflow.WorkflowManager) (*workflow.Workflow, error) {
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
		"Housekeeping (Direct)",
		"Performs routine housekeeping tasks like cleaning up temporary files and logs (direct trigger version)",
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

// triggerDirectWorkflowViaAPI triggers a workflow execution via the API
func triggerDirectWorkflowViaAPI(workflowID int, parameters json.RawMessage) error {
	// Create the request payload
	payload := map[string]interface{}{
		"workflow_id": workflowID,
		"parameters":  parameters,
		// No scheduled_at means it will run immediately
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(
		"POST",
		"http://localhost:3000/api/workflows/schedule",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your-token-here") // Replace with a valid token

	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	// Parse the response
	var execution struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&execution); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Workflow execution created with ID: %d\n", execution.ID)
	return nil
}

// connectToDB connects to the PostgreSQL database
func connectToDB() (*sql.DB, error) {
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
