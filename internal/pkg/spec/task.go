package spec

import (
	"encoding/json"
	"fmt"
	"github.com/glesica/flowork/internal/app/options"
	"github.com/glesica/flowork/internal/pkg/files"
	"io"
	"os"
)

type Task struct {
	// Name is a human-readable name for the task, to be used in UI and
	// logs as a quick way to reference a specific task. For example:
	// "parse", "train", "load data".
	Name string `json:"name" toml:"name"`

	// Desc is a human-readable description of the task, intended to be
	// included in UIs and documentation. This should generally be about
	// one sentence.
	Desc string `json:"desc" toml:"desc"`

	// The command to run as an array of strings equivalent
	// to an argv array, including the executable path.
	//
	// Examples:
	//   - []string{"ls", "-l", "/usr/bin"}
	Cmd []string `json:"cmd" toml:"cmd"`

	// The Docker image that the command will run in. The
	// working directory will be set automatically and mounted
	// at the location specified by WorkDir.
	//
	// Examples:
	//   - "debian:bookworm-slim"
	Image string `json:"image" toml:"image"`

	// The location within the container to mount the working directory
	// and from which the command will be run.
	//
	// Examples:
	//   - "/work"
	WorkDir string `json:"workdir" toml:"workdir"`

	// Inputs is a list of files that must exist, relative to the
	// working directory, in order for the task to run.
	// For now, these must be bare file names as they will only be
	// copied directly into the working directory. In the future,
	// full paths relative to the working directory will be supported.
	// TODO: Support full paths for inputs
	Inputs []files.Path `json:"inputs" toml:"inputs"`

	// Outputs is a list of files that are guaranteed to exist, relative
	// to the working directory, after the task has completed.
	Outputs []files.Path `json:"outputs" toml:"outputs"`

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
