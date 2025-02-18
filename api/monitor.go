package api

import (
	"encoding/json"
	"net/http"

	"github.com/r2unit/talos-core/controller"
	"github.com/r2unit/talos-core/queue"
)

func WorkersHandler(w http.ResponseWriter, r *http.Request) {
	workers := controller.GetWorkers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	tasks := queue.GetPendingTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
