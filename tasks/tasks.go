// core/tasks/tasks.go
package tasks

import (
	"errors"

	"github.com/r2unit/talos-core/queue" // replace with your actual module path
)

// TaskCreator is a function type that creates a queue.Task from given parameters.
type TaskCreator func(params map[string]string) (queue.Task, error)

// registry maps task types to their corresponding creator functions.
var registry = make(map[string]TaskCreator)

// RegisterTask registers a TaskCreator for a specific task type.
func RegisterTask(taskType string, creator TaskCreator) {
	registry[taskType] = creator
}

// GetTaskCreator retrieves a TaskCreator for the given task type.
func GetTaskCreator(taskType string) (TaskCreator, error) {
	creator, exists := registry[taskType]
	if !exists {
		return nil, errors.New("unknown task type")
	}
	return creator, nil
}
