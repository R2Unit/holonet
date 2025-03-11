package pgqueue

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/r2unit/holonet/queue"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func EnqueueTask(task queue.Task) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	argsJSON, err := json.Marshal(task.Args)
	if err != nil {
		return err
	}
	filesJSON, err := json.Marshal(task.Files)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO tasks(task_id, command, args, files, reporter, hosts, task_template, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', NOW(), NOW())
	`
	_, err = db.Exec(query, task.ID, task.Command, argsJSON, filesJSON, task.Reporter, task.Hosts, task.TaskTemplate)
	return err
}

func DequeueTask() (*queue.Task, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT task_id, command, args, files, reporter, hosts, task_template
		FROM tasks
		WHERE status = 'pending'
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`
	row := tx.QueryRow(query)
	var taskID, command, reporter, hosts, taskTemplate string
	var argsJSON, filesJSON []byte
	err = row.Scan(&taskID, &command, &argsJSON, &filesJSON, &reporter, &hosts, &taskTemplate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	_, err = tx.Exec(`UPDATE tasks SET status = 'running', updated_at = NOW() WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var args []string
	var files map[string]string
	if err := json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(filesJSON, &files); err != nil {
		return nil, err
	}

	task := &queue.Task{
		ID:           taskID,
		Command:      command,
		Args:         args,
		Files:        files,
		Reporter:     reporter,
		Hosts:        hosts,
		TaskTemplate: taskTemplate,
	}
	return task, nil
}

func UpdateTaskStatus(taskID, status string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := db.Exec(`UPDATE tasks SET status = $1, updated_at = NOW() WHERE task_id = $2`, status, taskID)
	return err
}
