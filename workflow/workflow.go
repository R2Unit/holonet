package workflow

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/holonet/core/logger"
)

type WorkflowStatus string

const (
	StatusDraft    WorkflowStatus = "draft"
	StatusActive   WorkflowStatus = "active"
	StatusInactive WorkflowStatus = "inactive"
	StatusArchived WorkflowStatus = "archived"
)

type ExecutionStatus string

const (
	ExecutionPending   ExecutionStatus = "pending"
	ExecutionRunning   ExecutionStatus = "running"
	ExecutionCompleted ExecutionStatus = "completed"
	ExecutionFailed    ExecutionStatus = "failed"
	ExecutionCancelled ExecutionStatus = "cancelled"
)

type Workflow struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Code        string         `json:"code"`
	Status      WorkflowStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type WorkflowExecution struct {
	ID           int             `json:"id"`
	WorkflowID   int             `json:"workflow_id"`
	Status       ExecutionStatus `json:"status"`
	Parameters   json.RawMessage `json:"parameters"`
	Result       json.RawMessage `json:"result"`
	ErrorMessage string          `json:"error_message"`
	ScheduledAt  time.Time       `json:"scheduled_at"`
	StartedAt    time.Time       `json:"started_at"`
	CompletedAt  time.Time       `json:"completed_at"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type WorkflowManager struct {
	db *sql.DB
}

func NewWorkflowManager(db *sql.DB) *WorkflowManager {
	return &WorkflowManager{db: db}
}

func (wm *WorkflowManager) CreateWorkflow(name, description, code string) (*Workflow, error) {
	workflow := &Workflow{
		Name:        name,
		Description: description,
		Code:        code,
		Status:      StatusDraft,
	}

	query := `
		INSERT INTO workflows (name, description, code, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := wm.db.QueryRow(
		query,
		workflow.Name,
		workflow.Description,
		workflow.Code,
		workflow.Status,
	).Scan(&workflow.ID, &workflow.CreatedAt, &workflow.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	return workflow, nil
}

func (wm *WorkflowManager) GetWorkflow(id int) (*Workflow, error) {
	query := `
		SELECT id, name, description, code, status, created_at, updated_at
		FROM workflows
		WHERE id = $1
	`

	workflow := &Workflow{}
	err := wm.db.QueryRow(query, id).Scan(
		&workflow.ID,
		&workflow.Name,
		&workflow.Description,
		&workflow.Code,
		&workflow.Status,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workflow not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	return workflow, nil
}

func (wm *WorkflowManager) UpdateWorkflow(workflow *Workflow) error {
	query := `
		UPDATE workflows
		SET name = $1, description = $2, code = $3, status = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	err := wm.db.QueryRow(
		query,
		workflow.Name,
		workflow.Description,
		workflow.Code,
		workflow.Status,
		workflow.ID,
	).Scan(&workflow.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	return nil
}

func (wm *WorkflowManager) ListWorkflows() ([]*Workflow, error) {
	query := `
		SELECT id, name, description, code, status, created_at, updated_at
		FROM workflows
		ORDER BY id
	`

	rows, err := wm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}
	defer rows.Close()

	workflows := []*Workflow{}
	for rows.Next() {
		workflow := &Workflow{}
		err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&workflow.Code,
			&workflow.Status,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow: %w", err)
		}
		workflows = append(workflows, workflow)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workflows: %w", err)
	}

	return workflows, nil
}

func (wm *WorkflowManager) ScheduleWorkflow(workflowID int, parameters json.RawMessage, scheduledAt time.Time) (*WorkflowExecution, error) {
	workflow, err := wm.GetWorkflow(workflowID)
	if err != nil {
		return nil, err
	}

	if workflow.Status != StatusActive {
		return nil, errors.New("cannot schedule inactive workflow")
	}

	execution := &WorkflowExecution{
		WorkflowID:  workflowID,
		Status:      ExecutionPending,
		Parameters:  parameters,
		ScheduledAt: scheduledAt,
	}

	query := `
		INSERT INTO workflow_executions (workflow_id, status, parameters, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err = wm.db.QueryRow(
		query,
		execution.WorkflowID,
		execution.Status,
		execution.Parameters,
		execution.ScheduledAt,
	).Scan(&execution.ID, &execution.CreatedAt, &execution.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to schedule workflow: %w", err)
	}

	logger.Info("Scheduled workflow %d for execution at %s", workflowID, scheduledAt)
	return execution, nil
}

func (wm *WorkflowManager) GetPendingExecutions() ([]*WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, status, parameters, result, error_message, scheduled_at, started_at, completed_at, created_at, updated_at
		FROM workflow_executions
		WHERE status = $1 AND scheduled_at <= NOW()
		ORDER BY scheduled_at
	`

	rows, err := wm.db.Query(query, ExecutionPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending executions: %w", err)
	}
	defer rows.Close()

	executions := []*WorkflowExecution{}
	for rows.Next() {
		execution := &WorkflowExecution{}
		err := rows.Scan(
			&execution.ID,
			&execution.WorkflowID,
			&execution.Status,
			&execution.Parameters,
			&execution.Result,
			&execution.ErrorMessage,
			&execution.ScheduledAt,
			&execution.StartedAt,
			&execution.CompletedAt,
			&execution.CreatedAt,
			&execution.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}
		executions = append(executions, execution)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating executions: %w", err)
	}

	return executions, nil
}

func (wm *WorkflowManager) UpdateExecution(execution *WorkflowExecution) error {
	query := `
		UPDATE workflow_executions
		SET status = $1, result = $2, error_message = $3, started_at = $4, completed_at = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := wm.db.QueryRow(
		query,
		execution.Status,
		execution.Result,
		execution.ErrorMessage,
		execution.StartedAt,
		execution.CompletedAt,
		execution.ID,
	).Scan(&execution.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update execution: %w", err)
	}

	return nil
}

func Init(db *sql.DB) error {
	logger.Info("Initializing workflow system")
	logger.Info("Workflow system initialized successfully")
	return nil
}
