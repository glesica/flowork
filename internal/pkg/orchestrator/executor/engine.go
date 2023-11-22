package executor

import (
	"fmt"
	"log/slog"

	"github.com/glesica/flowork/internal/pkg/orchestrator"
	"github.com/glesica/flowork/internal/pkg/task"
)

// An Engine arranges for a job to be run using its configured runner.
// This is a helper for the rest of the runner infrastructure to
// separate running the job from pooling and retries and such.
type Engine func(job *orchestrator.Job) error

// NoopEngine is the default, it does nothing. It is used for testing.
func NoopEngine(job *orchestrator.Job) error {
	slog.Info("executing job", "engine", "noop", "job", job.Id)
	return nil
}

// SimpleEngine will generally be used to run workflows in practice.
func SimpleEngine(job *orchestrator.Job) error {
	slog.Info("executing job", "engine", "simple", "job", job.Id)

	vol, err := job.Runner.CreateVolume(0)
	if err != nil {
		return fmt.Errorf("failed to create volume %w", err)
	}

	slog.Debug("created job volume", "engine", "simple", "job", job.Id, "volume", vol)

	defer func() {
		// TODO: Don't delete in debug mode
		err = job.Runner.DeleteVolume(vol)
		if err != nil {
			slog.Error("failed to delete volume", "engine", "simple", "error", err, "job", job.Id, "volume", vol)
		} else {
			slog.Debug("deleted job volume", "engine", "simple", "job", job.Id, "volume", vol)
		}
	}()

	firstTask := job.Tasks[0]
	for _, input := range firstTask.Inputs {
		err := job.Runner.AddFile(job.InPath, vol, input.File())
		if err != nil {
			return fmt.Errorf("simple engine: failed to copy %s to volume as %s: %w", job.InPath, input, err)
		} else {
			slog.Debug("copied input to volume", "engine", "simple", "job", job.Id, "volume", vol, "input", job.InPath)
		}
	}

	err = task.RunAll(job.Runner, job.Tasks, vol)
	if err != nil {
		return fmt.Errorf("simple engine: failed to run tasks: %w", err)
	}

	slog.Debug("finished running tasks", "engine", "simple", "job", job.Id, "volume", vol)

	lastTask := job.Tasks[len(job.Tasks)-1]
	for _, output := range lastTask.Task.Outputs {
		dest := job.OutDir.SubDir(lastTask.ID)

		err := job.Runner.ExtractFile(output, vol, dest)
		if err != nil {
			return fmt.Errorf("simple engine: failed to extract %s from volume as %s: %w", output, dest, err)
		}
	}

	slog.Debug("finished copying outputs", "engine", "simple", "job", job.Id, "volume", vol)

	return nil
}
