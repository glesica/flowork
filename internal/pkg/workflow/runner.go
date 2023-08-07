package workflow

import (
	"fmt"
	"github.com/glesica/flowork/internal/pkg/task"
)

func Run(w Instance, r task.Runner) error {
	for _, t := range w.Tasks {
		err := r.Run(t)
		if err != nil {
			return fmt.Errorf("failed to run workflow, task failure: %w", err)
		}
	}

	return nil
}
