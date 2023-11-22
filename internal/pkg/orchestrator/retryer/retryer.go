package retryer

import (
	"fmt"
	"log/slog"

	"github.com/glesica/flowork/internal/pkg/option"
	"github.com/glesica/flowork/internal/pkg/orchestrator"
)

// A worker handles failed jobs and decides whether they will
// be retried. The consumer never actually has access to
// the struct, it is just used internally to hold state.
type worker struct {
	errorQueue <-chan *orchestrator.Job
	retryQueue chan<- *orchestrator.Job

	maxRetries  int
	maxFailures int
}

// Start creates a new retry worker. It will read jobs that
// have errored from the given queue and determine whether they
// should be retried. If a job is to be retried, it will be sent
// on the queue returned by Retries().
//
// To properly shut down the retryer, close the error queue. The
// retry queue will be closed upon shutdown, so do not close it
// from the outside.
func Start(errorQueue <-chan *orchestrator.Job, retryQueue chan<- *orchestrator.Job, opts ...option.Func[*worker]) error {
	p := &worker{
		errorQueue: errorQueue,
		retryQueue: retryQueue,
	}

	if p.errorQueue == nil {
		return fmt.Errorf("an error queue channel is required")
	}

	err := option.Apply(p, opts...)
	if err != nil {
		return fmt.Errorf("failed to apply options to worker: %w", err)
	}

	go p.run()

	return nil
}

func (p *worker) run() {
	defer close(p.retryQueue)

	slog.Debug("retryer starting")

	failureCount := 0

	for {
		job, more := <-p.errorQueue
		if !more {
			break
		}

		if job.Attempts > p.maxRetries {
			// Job has already been tried the maximum number of times,
			// so we abandon it and consider it a failure
			failureCount++
			slog.Info("job failed", "id", job.Id, "inpath", job.InPath)
			continue
		}

		if failureCount >= p.maxFailures {
			// We surpassed our max failures, so we prepare to shut
			// down by no longer retrying jobs
			continue
		}

		job.Err = nil
		p.retryQueue <- job
	}

	slog.Debug("retryer finished")
}

// WithMaxRetries sets the maximum number of times a failed job
// will be retried before it becomes a failure.
func WithMaxRetries(value int) option.Func[*worker] {
	return func(p *worker) error {
		if value < 0 {
			return fmt.Errorf("max retries value must be non-negative: %d", value)
		}
		p.maxRetries = value
		return nil
	}
}

// WithMaxFailures sets the maximum number of failures (jobs that
// have run out of retries) that may occur before no more jobs
// will be retried, giving the rest of the system a chance to work
// through its queued jobs (whether they error or not) and shut down.
func WithMaxFailures(value int) option.Func[*worker] {
	return func(p *worker) error {
		if value < 0 {
			return fmt.Errorf("max failures value must be non-negative: %d", value)
		}
		p.maxFailures = value
		return nil
	}
}
