package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/holonet/core/logger"
)

type Executor struct {
	manager *WorkflowManager
}

func NewExecutor(manager *WorkflowManager) *Executor {
	return &Executor{
		manager: manager,
	}
}

func (e *Executor) StartExecutionLoop(ctx context.Context) {
	logger.Info("Starting workflow execution loop")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping workflow execution loop")
			return
		case <-ticker.C:
			if err := e.processWorkflows(); err != nil {
				logger.Error("Error processing workflows: %v", err)
			}
		}
	}
}

func (e *Executor) processWorkflows() error {
	executions, err := e.manager.GetPendingExecutions()
	if err != nil {
		return fmt.Errorf("failed to get pending executions: %w", err)
	}

	for _, execution := range executions {
		go e.executeWorkflow(execution)
	}

	return nil
}

func (e *Executor) executeWorkflow(execution *WorkflowExecution) {
	logger.Info("Executing workflow %d (execution %d)", execution.WorkflowID, execution.ID)
	execution.Status = ExecutionRunning
	execution.StartedAt = time.Now()
	if err := e.manager.UpdateExecution(execution); err != nil {
		logger.Error("Failed to update execution status to running: %v", err)
		return
	}

	workflow, err := e.manager.GetWorkflow(execution.WorkflowID)
	if err != nil {
		logger.Error("Failed to get workflow %d: %v", execution.WorkflowID, err)
		e.markExecutionFailed(execution, fmt.Sprintf("Failed to get workflow: %v", err))
		return
	}

	result, err := e.runWorkflowCode(workflow, execution.Parameters)
	if err != nil {
		logger.Error("Failed to execute workflow %d: %v", execution.WorkflowID, err)
		e.markExecutionFailed(execution, fmt.Sprintf("Execution error: %v", err))
		return
	}

	execution.Status = ExecutionCompleted
	execution.Result = result
	execution.CompletedAt = time.Now()
	if err := e.manager.UpdateExecution(execution); err != nil {
		logger.Error("Failed to update execution status to completed: %v", err)
		return
	}

	logger.Info("Workflow %d (execution %d) completed successfully", execution.WorkflowID, execution.ID)
}

func (e *Executor) markExecutionFailed(execution *WorkflowExecution, errorMessage string) {
	execution.Status = ExecutionFailed
	execution.ErrorMessage = errorMessage
	execution.CompletedAt = time.Now()
	if err := e.manager.UpdateExecution(execution); err != nil {
		logger.Error("Failed to update execution status to failed: %v", err)
	}
}

func (e *Executor) runWorkflowCode(workflow *Workflow, parameters json.RawMessage) (json.RawMessage, error) {
	logger.Info("Simulating execution of workflow %d: %s", workflow.ID, workflow.Name)
	time.Sleep(2 * time.Second)
	result := map[string]interface{}{
		"workflow_name": workflow.Name,
		"executed_at":   time.Now(),
		"parameters":    parameters,
		"result":        "Workflow executed successfully",
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	return resultJSON, nil
}
