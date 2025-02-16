// core/controller/controller.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/r2unit/talos-core/queue" // adjust module path as needed
	"github.com/r2unit/talos-core/web"   // uses the shared WebSocket code
)

const validToken = "secret123" // Shared token for authentication

// WSHandler upgrades the HTTP connection to a WebSocket using core/web,
// then continuously sends tasks (as JSON) from the core's task queue,
// and listens for status updates from the worker.
func WSHandler(w http.ResponseWriter, r *http.Request) {
	// Capture the worker name from query parameters.
	workerName := r.URL.Query().Get("name")

	conn, err := web.Upgrade(w, r, validToken)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	// Ensure the connection is closed and log when the worker disconnects.
	defer func() {
		conn.Close()
		log.Printf("Worker '%s' disconnected", workerName)
	}()

	// Start a goroutine to read status messages from the worker.
	go func() {
		for {
			msg, err := web.ReadMessage(conn)
			if err != nil {
				log.Println("Error reading status message:", err)
				return
			}
			log.Printf("Status update from worker '%s': %s", workerName, msg)
		}
	}()

	// Continuously send tasks to the worker.
	for {
		select {
		case task := <-queue.TaskQueue:
			msg, err := json.Marshal(task)
			if err != nil {
				log.Println("JSON Marshal error:", err)
				continue
			}
			if err := web.WriteMessage(conn, string(msg)); err != nil {
				log.Println("WriteMessage error:", err)
				return
			}
			log.Printf("Sent task %s to worker '%s'", task.ID, workerName)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
