package task

import (
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/spec"
)

type Instance struct {
	spec.Task

	// ID is the unique identifier for the task instance.
	ID string `json:"id"`
}

// TODO: Can we take input/output handling away from the runners?

func NewInstance(t spec.Task) Instance {
	return Instance{
		Task: t,
		ID:   id.New(),
	}
}
