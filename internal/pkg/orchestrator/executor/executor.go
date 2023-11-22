package executor

import (
	"fmt"
	"log/slog"

	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/option"
	"github.com/glesica/flowork/internal/pkg/orchestrator"
)

type worker struct {
	inputQueue   <-chan files.Path
	errorQueue   chan<- *orchestrator.Job
	successQueue chan<- *orchestrator.Job

	createJob JobCreator
	engine    Engine
}

func Start(inputQueue <-chan files.Path, options ...option.Func[*worker]) error {
	w := &worker{
		inputQueue: inputQueue,
	}

	err := option.Apply(w, options...)
	if err != nil {
		return fmt.Errorf("failed to start executor: %w", err)
	}

	// The default job creator includes no tasks, so the default
	// job will do nothing at all, suitable for testing workflow
	// structure.
	if w.createJob == nil {
		w.createJob = MakeJobCreator(nil)
	}

	// The default engine does nothing, so the default is suitable
	// for testing and won't screw anything up or incur costs.
	if w.engine == nil {
		w.engine = NoopEngine
	}

	if w.errorQueue == nil {
		w.errorQueue = make(chan *orchestrator.Job, 10)
	}

	if w.successQueue == nil {
		w.successQueue = make(chan *orchestrator.Job, 10)
	}

	go w.run()

	return nil
}

func (w *worker) run() {
	defer close(w.errorQueue)
	defer close(w.successQueue)

	for {
		input, more := <-w.inputQueue
		if !more {
			break
		}

		job, err := w.createJob(input)
		if err != nil {
			slog.Error("encountered error creating job, skipping", "error", err, "input", input)
			continue
		}

		err = w.engine(job)
		if err != nil {
			job.Err = err
			w.errorQueue <- job
		} else {
			w.successQueue <- job
		}
	}
}

func WithJobCreator(createJob JobCreator) option.Func[*worker] {
	return func(w *worker) error {
		w.createJob = createJob
		return nil
	}
}

func WithEngine(engine Engine) option.Func[*worker] {
	return func(w *worker) error {
		w.engine = engine
		return nil
	}
}

func WithErrorQueue(errorQueue chan<- *orchestrator.Job) option.Func[*worker] {
	return func(w *worker) error {
		w.errorQueue = errorQueue
		return nil
	}
}

func WithSuccessQueue(successQueue chan<- *orchestrator.Job) option.Func[*worker] {
	return func(w *worker) error {
		w.successQueue = successQueue
		return nil
	}
}
