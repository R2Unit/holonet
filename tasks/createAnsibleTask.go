package tasks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/r2unit/talos-core/queue"
)

// createAnsibleTask creates a new Task for running an ansible-playbook.
// It expects a parameter "hosts" (a comma-separated list of hosts).
// The inventory and playbook filenames are hardcoded to "inventory.ini" and "site.yml"
// (stored in core/template). The playbook file should contain the placeholder {{HOSTS}},
// which is replaced with the provided hosts.
func createAnsibleTask(params map[string]string) (queue.Task, error) {
	hosts, ok := params["hosts"]
	if !ok {
		return queue.Task{}, errors.New("missing 'hosts' parameter")
	}

	inventoryFile := "inventory.ini"
	playbookFile := "site.yml"
	templateDir := "core/template"

	inventoryPath := filepath.Join(templateDir, inventoryFile)
	inventoryContent, err := ioutil.ReadFile(inventoryPath)
	if err != nil {
		return queue.Task{}, fmt.Errorf("failed to read inventory file: %w", err)
	}

	playbookPath := filepath.Join(templateDir, playbookFile)
	playbookBytes, err := ioutil.ReadFile(playbookPath)
	if err != nil {
		return queue.Task{}, fmt.Errorf("failed to read playbook file: %w", err)
	}

	playbookContent := strings.Replace(string(playbookBytes), "{{HOSTS}}", hosts, -1)

	task := queue.Task{
		ID:      fmt.Sprintf("ansible-%d", time.Now().UnixNano()),
		Command: "ansible-playbook",
		Args:    []string{"-i", inventoryFile, playbookFile},
		Files: map[string]string{
			inventoryFile: string(inventoryContent),
			playbookFile:  playbookContent,
		},
	}
	return task, nil
}

func init() {
	RegisterTask("ansible", createAnsibleTask)
}
