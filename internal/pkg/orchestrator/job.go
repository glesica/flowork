package orchestrator

import (
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/task"
)

// A Job represents a single, parallel, self-contained unit of
// work. It also tracks its own state as it moves through the
// execution machinery.
type Job struct {
	Id       string
	Runner   task.Runner
	Tasks    []*task.Instance
	InPath   files.Path
	OutDir   files.Dir
	Err      error
	Attempts int
}
