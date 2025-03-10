package api

import (
	"encoding/json"
	"net/http"

	"github.com/r2unit/holonet/queue"
	"github.com/r2unit/holonet/tasks"
)

type TaskRequest struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	creator, err := tasks.GetTaskCreator(req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newTask, err := creator(req.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	queue.Enqueue(newTask)
	resp := APIResponse{Status: "ok", Message: "Task enqueued"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
