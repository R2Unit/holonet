package main

import (
	"fmt"
	"time"
)

// HousekeepingWorkflow is an example workflow that performs basic housekeeping tasks
// This code would be stored in the workflow's "code" field in the database
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

// This is the workflow code that would be stored in the database
const HousekeepingWorkflowCode = `
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
