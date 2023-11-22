package workflow

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/option"
	"github.com/glesica/flowork/internal/pkg/spec"
)

type Instance struct {
	spec.Workflow
	ID    string       `json:"id"`
	Tasks spec.TaskSet `json:"tasks"`

	// TODO: Where to store outputs? Should that go here?
	// We could add workflow-level options on to the instance

	// CaptureDir is the directory where data from this workflow run
	// will be written, such as logs, stdout, stderr, and so on.
	// The default value is the current working directory with the
	// instance ID appended to it.
	CaptureDir files.Dir
}

func NewInstance(w spec.Workflow, opts ...option.Func[*Instance]) (*Instance, error) {
	instance := &Instance{
		Workflow: w,
		ID:       id.New(),
		Tasks:    w.Tasks[:],
	}

	err := option.Apply(instance, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to apply all options: %w", err)
	}

	if instance.CaptureDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}

		err = WithWorkDir(files.Dir(cwd))(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to apply workdir option: %w", err)
		}
	}

	slog.Info("created new workflow instance", "instance", instance)

	return instance, nil
}

func WithWorkDir(d files.Dir) option.Func[*Instance] {
	return func(instance *Instance) error {
		instance.CaptureDir = d.SubDir(instance.ID)
		return nil
	}
}
