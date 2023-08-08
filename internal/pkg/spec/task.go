package spec

import (
	"encoding/json"
	"fmt"
	"github.com/glesica/flowork/internal/app/options"
	"io"
	"os"
)

type Task struct {
	// The command to run as an array of strings equivalent
	// to an argv array, including the executable path.
	//
	// Examples:
	//   - []string{"ls", "-l", "/usr/bin"}
	Cmd []string `json:"cmd"`

	// The Docker image that the command will run in. The
	// working directory will be set automatically and mounted
	// at the location specified by WorkDir.
	//
	// Examples:
	//   - "debian:bookworm-slim"
	Image string `json:"image"`

	// The location within the container to mount the working directory
	// and from which the command will be run.
	//
	// Examples:
	//   - "/work"
	WorkDir string `json:"workdir"`

	// Desc is a human-readable description of the task, intended to be
	// included in UIs and documentation.
	Desc string `json:"desc"`

	// Inputs - files that must exist in the working directory before
	// the command can run, optional
	// TODO: Inputs []string `json:"inputs"`

	// TODO: We could also delete everything but the outputs for efficiency

	// Outputs - files that must exist in the working directory after
	// the command has run in order to consider the task a success,
	// optional
	// TODO: Outputs []string `json:"outputs"`

	// DiskSpaceGB indicates the required amount of disk space
	// available on the volume where the working directory is
	// located. The actual amount may be larger, but it will
	// not be smaller.
	// TODO: DiskSpaceGB uint `json:"disk-space-gb"`

	// MemoryGB indicates the minimum number of GB of RAM the
	// task requires. The actual amount may be larger, but it
	// will not be smaller.
	// TODO: MemoryGB uint `json:"memory-gb"`
}

func LoadTaskPath(p string) (Task, error) {
	f, err := os.Open(p)
	if err != nil {
		return Task{}, fmt.Errorf("failed to open file (%s): %w", p, err)
	}
	defer func() { _ = f.Close() }()

	return LoadTask(f)
}

func LoadTask(data io.Reader) (Task, error) {
	c := Task{}

	raw, err := io.ReadAll(data)
	if err != nil {
		return c, fmt.Errorf("failed to load task from JSON: %w", err)
	}

	err = json.Unmarshal(raw, &c)
	if err != nil {
		return c, fmt.Errorf("failed to parse task from JSON: %w", err)
	}

	return c, nil
}

// GetWorkDir returns the working directory path for the task
// (see WorkDir) but returns the default value if the field is
// blank (the empty string).
func (t Task) GetWorkDir() string {
	if t.WorkDir == "" {
		return options.DefaultTaskWorkDir
	}

	return t.WorkDir
}

// TaskSet is a collection of tasks that can be assigned to
// a workflow.
type TaskSet []Task
