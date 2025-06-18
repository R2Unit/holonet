package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/holonet/core/logger"
	"github.com/holonet/core/workflow"
)

var workflowManager *workflow.WorkflowManager

func SetWorkflowManager(manager *workflow.WorkflowManager) {
	workflowManager = manager
}

func RegisterWorkflowRoutes() {
	http.HandleFunc("/api/workflows", tokenAuthMiddleware(handleWorkflows))
	http.HandleFunc("/api/workflows/", tokenAuthMiddleware(handleWorkflowByID))
	http.HandleFunc("/api/workflows/schedule", tokenAuthMiddleware(handleScheduleWorkflow))
}

func handleWorkflows(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWorkflows(w, r)
	case http.MethodPost:
		createWorkflow(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWorkflowByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/workflows/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getWorkflow(w, r, id)
	case http.MethodPut:
		updateWorkflow(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleScheduleWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		WorkflowID  int             `json:"workflow_id"`
		Parameters  json.RawMessage `json:"parameters"`
		ScheduledAt string          `json:"scheduled_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var scheduledAt time.Time
	var err error
	if request.ScheduledAt != "" {
		scheduledAt, err = time.Parse(time.RFC3339, request.ScheduledAt)
		if err != nil {
			http.Error(w, "Invalid scheduled_at format. Use RFC3339 format (e.g., 2025-01-01T12:00:00Z)", http.StatusBadRequest)
			return
		}
	} else {
		scheduledAt = time.Now()
	}

	execution, err := workflowManager.ScheduleWorkflow(request.WorkflowID, request.Parameters, scheduledAt)
	if err != nil {
		logger.Error("Failed to schedule workflow: %v", err)
		http.Error(w, "Failed to schedule workflow: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(execution)
}

func getWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := workflowManager.ListWorkflows()
	if err != nil {
		logger.Error("Failed to list workflows: %v", err)
		http.Error(w, "Failed to list workflows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

func createWorkflow(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Code        string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workflow, err := workflowManager.CreateWorkflow(request.Name, request.Description, request.Code)
	if err != nil {
		logger.Error("Failed to create workflow: %v", err)
		http.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workflow)
}

func getWorkflow(w http.ResponseWriter, r *http.Request, id int) {
	wf, err := workflowManager.GetWorkflow(id)
	if err != nil {
		logger.Error("Failed to get workflow: %v", err)
		http.Error(w, "Failed to get workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wf)
}

func updateWorkflow(w http.ResponseWriter, r *http.Request, id int) {
	var request struct {
		Name        string                  `json:"name"`
		Description string                  `json:"description"`
		Code        string                  `json:"code"`
		Status      workflow.WorkflowStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wf, err := workflowManager.GetWorkflow(id)
	if err != nil {
		logger.Error("Failed to get workflow: %v", err)
		http.Error(w, "Failed to get workflow", http.StatusInternalServerError)
		return
	}

	wf.Name = request.Name
	wf.Description = request.Description
	wf.Code = request.Code
	wf.Status = request.Status

	if err := workflowManager.UpdateWorkflow(wf); err != nil {
		logger.Error("Failed to update workflow: %v", err)
		http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wf)
}
