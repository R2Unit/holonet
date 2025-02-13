// core/tasks/createAnsibleTask.go
package tasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/r2unit/talos-core/queue"
)

func createAnsibleTask(params map[string]string) (queue.Task, error) {
	inventory, ok := params["inventory"]
	if !ok {
		return queue.Task{}, errors.New("missing 'inventory' parameter")
	}
	playbook, ok := params["playbook"]
	if !ok {
		return queue.Task{}, errors.New("missing 'playbook' parameter")
	}
	task := queue.Task{
		ID:      fmt.Sprintf("ansible-%d", time.Now().UnixNano()),
		Command: "ansible-playbook",
		Args:    []string{"-i", inventory, playbook},
	}
	return task, nil
}

// Automatically register the "ansible" task creator when this package is initialized.
func init() {
	RegisterTask("ansible", createAnsibleTask)
}
