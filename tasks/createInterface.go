// core/tasks/createAnsibleTask.go
package tasks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/r2unit/holonet/queue" // Replace with your actual module path
)

// createAnsibleTask creates a new Task for running an ansible-playbook.
// It expects parameters:
//   - "hosts": a comma-separated list of target hosts,
//   - "reporter": who requested this task,
//   - "interface": the network interface to use (optional, default is "default").
//
// The inventory and playbook filenames are hardcoded to "inventory.ini" and "site.yml"
// (stored in core/template). The playbook file should contain placeholders
// like {{HOSTS}} and {{INTERFACE}} which will be replaced with the provided values.
// If a requirements.yml file exists in core/template, it is also attached.
func createInterface(params map[string]string) (queue.Task, error) {
	hosts, ok := params["hosts"]
	if !ok {
		return queue.Task{}, errors.New("missing 'hosts' parameter")
	}
	reporter, ok := params["reporter"]
	if !ok {
		return queue.Task{}, errors.New("missing 'reporter' parameter")
	}
	iface, ok := params["interface"]
	if !ok {
		// If not provided, use a default value.
		iface = "default"
	}

	// Hardcoded filenames and template directory.
	templateDir := "core/template"
	inventoryFile := "inventory.ini"
	playbookFile := "site.yml"
	reqFile := "requirements.yml"

	// Read the inventory file.
	invPath := filepath.Join(templateDir, inventoryFile)
	invContent, err := ioutil.ReadFile(invPath)
	if err != nil {
		return queue.Task{}, fmt.Errorf("failed to read inventory file: %w", err)
	}

	// Read the playbook template.
	playbookPath := filepath.Join(templateDir, playbookFile)
	playbookBytes, err := ioutil.ReadFile(playbookPath)
	if err != nil {
		return queue.Task{}, fmt.Errorf("failed to read playbook file: %w", err)
	}

	// Replace placeholders in the playbook template.
	playbookContent := string(playbookBytes)
	playbookContent = strings.ReplaceAll(playbookContent, "{{HOSTS}}", hosts)
	playbookContent = strings.ReplaceAll(playbookContent, "{{INTERFACE}}", iface)

	// Optionally, read the requirements file if it exists.
	var reqContent string
	reqPath := filepath.Join(templateDir, reqFile)
	if data, err := ioutil.ReadFile(reqPath); err == nil {
		reqContent = string(data)
	} else {
		reqContent = ""
	}

	// Create the task.
	task := queue.Task{
		ID:      fmt.Sprintf("ansible-%d", time.Now().UnixNano()),
		Command: "ansible-playbook",
		Args:    []string{"-i", inventoryFile, playbookFile},
		Files: map[string]string{
			inventoryFile: string(invContent),
			playbookFile:  playbookContent,
		},
		Reporter:     reporter,
		Hosts:        hosts,
		TaskTemplate: "ansible",
	}

	// Attach requirements file if available.
	if reqContent != "" {
		task.Files[reqFile] = reqContent
	}

	return task, nil
}

func init() {
	RegisterTask("add_interface", createInterface)
}
