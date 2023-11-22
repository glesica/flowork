package task

import (
	"errors"
	"fmt"
	"io"

	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/option"
	"github.com/glesica/flowork/internal/pkg/spec"
)

// An Instance provides everything necessary for a runner to interact
// with a particular task in a structured way. It is mutable, and
// should only be submitted to a runner once, for this reason. Once
// it has passed through a runner, the workflow runtime will call
// Finalize(), after which it becomes read-only.
type Instance struct {
	spec.Task

	// ID is the unique identifier for the task instance.
	ID string `json:"id"`

	stdout io.WriteCloser
	stderr io.WriteCloser
}

func NewInstance(t spec.Task, opts ...option.Func[*Instance]) (*Instance, error) {
	instance := &Instance{
		Task: t,
		ID:   id.New(),
	}

	err := option.Apply(instance, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to apply all options: %w", err)
	}

	return instance, nil
}

func WithStdout(w io.WriteCloser) option.Func[*Instance] {
	return func(instance *Instance) error {
		instance.stdout = w
		return nil
	}
}

func WithStderr(w io.WriteCloser) option.Func[*Instance] {
	return func(instance *Instance) error {
		instance.stderr = w
		return nil
	}
}

func (t *Instance) Stdout() io.Writer {
	return t.stdout
}

func (t *Instance) Stderr() io.Writer {
	return t.stderr
}

func (t *Instance) Finalize() error {
	var errs []error

	errs = append(errs, t.stdout.Close())
	errs = append(errs, t.stderr.Close())

	return errors.Join(errs...)
}
