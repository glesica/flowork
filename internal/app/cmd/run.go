package cmd

import (
	"fmt"
	"github.com/glesica/flowork/internal/pkg/spec"
	"github.com/glesica/flowork/internal/pkg/task"
	"github.com/glesica/flowork/internal/pkg/workflow"
)

type RunOptions struct {
	Workflow string `help:"Path to workflow definition to execute"`
	Runner   string `help:"Task runner to use" enum:"docker" default:"docker"`
}

func (r *RunOptions) Run(global GlobalOptions) error {
	ws, err := spec.LoadWorkflowPath(r.Workflow)
	if err != nil {
		return fmt.Errorf("failed to load workflow (%s): %w", r.Workflow, err)
	}

	wi := workflow.NewInstance(ws)

	var runner task.Runner
	switch r.Runner {
	case "docker":
		runner = &task.DockerRunner{}
	}

	err = workflow.Run(wi, runner)
	if err != nil {
		return fmt.Errorf("failed to run workflow (%s): %w", r.Workflow, err)
	}

	return nil
}
