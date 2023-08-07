package workflow

import (
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/spec"
	"github.com/glesica/flowork/internal/pkg/task"
	"golang.org/x/exp/slog"
)

type Instance struct {
	spec.Workflow

	ID    string          `json:"id"`
	Tasks []task.Instance `json:"tasks"`
}

func NewInstance(w spec.Workflow) Instance {
	inst := Instance{
		Workflow: w,

		ID:    id.New(),
		Tasks: nil,
	}

	for _, t := range w.Tasks {
		inst.Tasks = append(inst.Tasks, task.NewInstance(t))
	}

	slog.Info("created new workflow instance", "instance", inst)

	return inst
}
