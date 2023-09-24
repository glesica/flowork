package workflow

import (
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/task"
)

type Job struct {
	Runner   task.Runner
	Tasks    []task.Instance
	InPaths  []files.Path
	OutDir   files.Dir
	Err      error
	Attempts int
}
