package recorder

import (
	"github.com/glesica/flowork/internal/pkg/option"
	"github.com/glesica/flowork/internal/pkg/orchestrator"
)

// A worker asynchronously captures the results of finished
// jobs using whatever method has been configured.
type worker struct {
	targets []Target
}

// Start runs the worker in the background, reading its input
// from the given channels. Once these channels have been closed,
// it will shut down.
func Start(successQueue <-chan *orchestrator.Job, failureQueue <-chan *orchestrator.Job, options ...option.Func[*worker]) {
	//
}

func (w *worker) run() {
	//
}

func WithTarget(t Target) option.Func[*worker] {
	return func(w *worker) error {
		w.targets = append(w.targets, t)
		return nil
	}
}
