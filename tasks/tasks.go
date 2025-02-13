package tasks

import (
	"errors"

	"github.com/r2unit/talos-core/queue"
)

type TaskCreator func(params map[string]string) (queue.Task, error)

var registry = make(map[string]TaskCreator)

func RegisterTask(taskType string, creator TaskCreator) {
	registry[taskType] = creator
}

func GetTaskCreator(taskType string) (TaskCreator, error) {
	creator, exists := registry[taskType]
	if !exists {
		return nil, errors.New("unknown task type")
	}
	return creator, nil
}
