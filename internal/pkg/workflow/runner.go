package workflow

import (
	"context"
	"errors"
	"fmt"
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/task"
	"golang.org/x/sync/semaphore"
	"log/slog"
	"sync"
)

func Run(w Instance, runner task.Runner, inFiles files.Iter, outDir files.Dir, concurrency int64) error {
	// TODO: Maybe add a task.Result type that bundles outputs and errors?

	maxAttempts := 3
	maxErrs := 1

	ctx, cancel := context.WithCancel(context.Background())

	sem := semaphore.NewWeighted(concurrency)

	var errs []error
	todo := make(chan *Job, concurrency+1)
	done := make(chan *Job, concurrency+1)

	processing := &sync.Mutex{}
	processing.Lock()

	// Start the job processor
	go func(jobQueue <-chan *Job, doneQueue chan<- *Job) {
		defer close(doneQueue)

		for {
			select {
			case _ = <-ctx.Done():
				return
			case jobNext, more := <-jobQueue:
				if !more {
					return
				}

				err := sem.Acquire(ctx, 1)
				if err != nil {
					return
				}

				go func(job *Job) {
					err := processJob(ctx, job)
					if err != nil {
						job.Err = err
					}

					job.Attempts++

					done <- job
				}(jobNext)
			}
		}
	}(todo, done)

	// Start the results processor
	go func(jobQueue chan<- *Job, doneQueue <-chan *Job) {
		defer processing.Unlock()

		for {
			select {
			case _ = <-ctx.Done():
				return
			case jobDone, more := <-done:
				if !more {
					return
				}

				sem.Release(1)

				if jobDone.Err != nil {
					errs = append(errs, jobDone.Err)
					if len(errs) > maxErrs {
						return
					}

					// Re-submit if we haven't exhausted the allowed attempts
					if jobDone.Attempts < maxAttempts {
						jobDone.Err = nil

						err := sem.Acquire(ctx, 1)
						if err != nil {
							return
						}

						jobQueue <- jobDone
					}
				}
			}
		}
	}(todo, done)

	// Submit one job for each input file
	for {
		inPath, ok := inFiles()
		if !ok {
			close(todo)
			break
		}

		var taskInstances []task.Instance
		for _, t := range w.Tasks {
			inst := task.NewInstance(t)
			taskInstances = append(taskInstances, inst)
		}

		jobNext := &Job{
			Tasks:   taskInstances,
			InPaths: []files.Path{inPath},
			OutDir:  outDir,
			Runner:  runner,
		}
		todo <- jobNext
	}

	// Wait for processing to complete, one way or another
	processing.Lock()
	cancel()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func processJob(ctx context.Context, job *Job) error {
	vol, err := job.Runner.CreateVolume(0)
	if err != nil {
		return fmt.Errorf("failed to create volume %w", err)
	}
	defer func() {
		// TODO: Don't delete in debug mode
		err = job.Runner.DeleteVolume(vol)
		if err != nil {
			slog.Error("failed to delete volume", "error", err, "volume", vol)
		}
	}()

	firstTask := job.Tasks[0]
	for _, input := range firstTask.Inputs {
		// TODO: Match a file to each of the inputs
		inPath := job.InPaths[0]

		err := job.Runner.AddFile(inPath, vol, input.File())
		if err != nil {
			return fmt.Errorf("failed to copy %s to volume as %s: %w", inPath, input, err)
		}
	}

	err = task.RunAll(job.Runner, job.Tasks, vol)
	if err != nil {
		return fmt.Errorf("failed to run tasks: %w", err)
	}

	lastTask := job.Tasks[len(job.Tasks)-1]
	for _, output := range lastTask.Task.Outputs {
		dest := job.OutDir.SubDir(lastTask.ID)

		err := job.Runner.ExtractFile(output, vol, dest)
		if err != nil {
			return fmt.Errorf("failed to extract %s from volume as %s: %w", output, dest, err)
		}
	}

	return nil
}
