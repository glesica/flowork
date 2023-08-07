package spec

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Workflow struct {
	// Name is the workflow name, used for UI purposes only.
	Name string `json:"name"`

	// Desc is a description of the workflow, intended for
	// UI and documentation purposes.
	Desc string `json:"desc"`

	// Tasks is the list of tasks to execute when the workflow
	// is run.
	Tasks TaskSet `json:"tasks"`
}

func LoadWorkflowPath(p string) (Workflow, error) {
	f, err := os.Open(p)
	if err != nil {
		return Workflow{}, fmt.Errorf("failed to open workflow file (%s): %w", p, err)
	}
	defer func() { _ = f.Close() }()

	return LoadWorkflow(f)
}

func LoadWorkflow(data io.Reader) (Workflow, error) {
	w := Workflow{}

	raw, err := io.ReadAll(data)
	if err != nil {
		return w, fmt.Errorf("failed to load workflow from file: %w", err)
	}

	err = json.Unmarshal(raw, &w)
	if err != nil {
		return w, fmt.Errorf("failed to parse workflow from JSON: %w", err)
	}

	return w, nil
}
