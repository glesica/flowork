package workflow

import (
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/spec"
	"log/slog"
)

type Instance struct {
	spec.Workflow
	ID    string       `json:"id"`
	Tasks spec.TaskSet `json:"tasks"`

	// TODO: Where to store outputs? Should that go here?
	// We could add workflow-level options on to the instance
}

func NewInstance(w spec.Workflow) Instance {
	inst := Instance{
		Workflow: w,
		ID:       id.New(),
		Tasks:    w.Tasks[:],
	}

	slog.Info("created new workflow instance", "instance", inst)

	return inst
}
