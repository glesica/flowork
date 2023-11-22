package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glesica/flowork/internal/app/options"
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/inputs"
	"github.com/glesica/flowork/internal/pkg/spec"
	"github.com/glesica/flowork/internal/pkg/task"
	"github.com/glesica/flowork/internal/pkg/workflow"
)

type RunOptions struct {
	Name        string    `help:"A human-readable name for this workflow run, will be used as a directory name"`
	Workflow    string    `help:"Path to workflow definition to execute" arg:""`
	Runner      string    `help:"Task runner to use" enum:"docker" default:"docker"`
	WorkDir     files.Dir `help:"Local working directory to use" default:"."`
	Input       files.Dir `help:"A directory to load inputs from"`
	Output      files.Dir `help:"A directory to save the outputs"`
	Concurrency int64     `help:"Max number of concurrent jobs (<1 means unlimited)" default:"1"`
}

func (o *RunOptions) setName() error {
	if o.Name == "" {
		y, m, d := time.Now().Date()
		idPart := id.New()

		o.Name = fmt.Sprintf("%d-%d-%d-%s", y, m, d, idPart)
	}

	return nil
}

func (o *RunOptions) setWorkDir() error {
	workDir := string(o.WorkDir)

	if strings.TrimSpace(workDir) == "." {
		// Assign the full working directory for good measure.
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		workDir = wd
	} else {
		// Turn the path we were given into an absolute path since
		// that makes everything easier (like dealing with Docker).
		wd, err := filepath.Abs(workDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute working directory: %w", err)
		}

		workDir = wd
	}

	o.WorkDir = files.Dir(workDir).SubDir(o.Name)

	return nil
}

func (o *RunOptions) setOutput() error {
	if o.Output == "" {
		o.Output = o.WorkDir.SubDir(options.OutputsDirName)
	}

	return nil
}

func (o *RunOptions) setConcurrency() error {
	if o.Concurrency < 1 {
		o.Concurrency = math.MaxInt
	}

	return nil
}

func Run(run *RunOptions, global GlobalOptions) error {
	var err error

	err = run.setName()
	if err != nil {
		return fmt.Errorf("failed to set run name: %w", err)
	}

	// setName must be called before setWorkDir
	err = run.setWorkDir()
	if err != nil {
		return fmt.Errorf("failed to set working directory: %w", err)
	}

	// setWorkDir must be called before setOutput
	err = run.setOutput()
	if err != nil {
		return fmt.Errorf("failed to set output location: %w", err)
	}

	err = run.setConcurrency()
	if err != nil {
		return fmt.Errorf("failed to set concurrency: %w", err)
	}

	ws, err := spec.LoadWorkflowPath(run.Workflow)
	if err != nil {
		return fmt.Errorf("failed to load workflow (%s): %w", run.Workflow, err)
	}

	wi := workflow.NewInstance(ws)

	var runner task.Runner
	switch run.Runner {
	case "docker":
		runner = &task.DockerRunner{
			Debug:   global.Debug,
			WorkDir: run.WorkDir,
			Store:   &files.Local{},
		}
	default:
		return fmt.Errorf("invalid runner (%s)", run.Runner)
	}

	in, err := inputs.Local(run.Input)
	if err != nil {
		return fmt.Errorf("failed to load inputs: %w", err)
	}

	err = workflow.Run(wi, runner, in, run.Output, run.Concurrency)
	if err != nil {
		return fmt.Errorf("failed to run workflow (%s): %w", run.Workflow, err)
	}

	return nil
}
