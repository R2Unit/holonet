package tasks

import (
	"errors"

	"github.com/r2unit/holonet/queue"
)

type TaskCreator func(params map[string]string) (queue.Task, error)

var registry = make(map[string]TaskCreator)

func RegisterTask(taskType string, creator TaskCreator) {
	registry[taskType] = creator
}

// TODO: Een betere task management voor als er undefined issues zijn.
func GetTaskCreator(taskType string) (TaskCreator, error) {
	creator, exists := registry[taskType]
	if !exists {
		return nil, errors.New("unknown task type")
	}
	return creator, nil
}
