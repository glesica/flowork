package executor

import (
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/orchestrator"
	"github.com/glesica/flowork/internal/pkg/spec"
	"github.com/glesica/flowork/internal/pkg/task"
)

type JobCreator func(inPath files.Path) (*orchestrator.Job, error)

func MakeJobCreator(tasks spec.TaskSet) JobCreator {
	return func(inPath files.Path) (*orchestrator.Job, error) {
		var taskInsts []*task.Instance
		for _, t := range tasks {
			inst, err := task.NewInstance(t)
			if err != nil {
				return nil, err
			}
			taskInsts = append(taskInsts, inst)
		}

		return &orchestrator.Job{
			Id:     id.New(),
			Tasks:  taskInsts,
			InPath: inPath,
		}, nil
	}
}
