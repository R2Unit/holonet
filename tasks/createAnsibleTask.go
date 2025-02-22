package tasks

import (
	"errors"
	"fmt"
	"io/ioutil"

	"path/filepath"
	"strings"
	"time"

	"github.com/r2unit/holonet-core/queue"
)

// createAnsibleTask creates a new Task for running an ansible-playbook.
// It expects parameters:
//   - "hosts": a comma-separated list of target hosts,
//   - "reporter": the identity of the requester.
//
// The inventory and playbook are read from core/template (hardcoded filenames).
// The playbook file should contain the placeholder {{HOSTS}} which will be replaced
// by the provided hosts. If a requirements.yml file exists in core/template, it is also attached.
func createAnsibleTask(params map[string]string) (queue.Task, error) {
	hosts, ok := params["hosts"]
	if !ok {
		return queue.Task{}, errors.New("missing 'hosts' parameter")
	}
	reporter, ok := params["reporter"]
	if !ok {
		return queue.Task{}, errors.New("missing 'reporter' parameter")
	}

	templateDir := "template"
	inventoryFile := "inventory.ini"
	playbookFile := "site.yml"
	reqFile := "requirements.yml"

	invPath := filepath.Join(templateDir, inventoryFile)
	invContent, err := ioutil.ReadFile(invPath)
	if err != nil {
		return queue.Task{}, fmt.Errorf("failed to read inventory file: %w", err)
	}

	playbookPath := filepath.Join(templateDir, playbookFile)
	playbookBytes, err := ioutil.ReadFile(playbookPath)
	if err != nil {
		return queue.Task{}, fmt.Errorf("failed to read playbook file: %w", err)
	}
	playbookContent := strings.Replace(string(playbookBytes), "{{HOSTS}}", hosts, -1)

	var reqContent string
	reqPath := filepath.Join(templateDir, reqFile)
	if data, err := ioutil.ReadFile(reqPath); err == nil {
		reqContent = string(data)
	} else {
		reqContent = ""
	}

	// Aanmaken van de task voor een Worker met alle aligned informatie
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

	// Controle voor de requirements van deze task
	if reqContent != "" {
		task.Files[reqFile] = reqContent
	}

	return task, nil
}

func init() {
	RegisterTask("ansible", createAnsibleTask)
}
