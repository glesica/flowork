package inputs

import (
	"log/slog"

	"github.com/glesica/flowork/internal/pkg/files"
)

// Iterator iterates over a collection of Paths
// and sends each over the channel it returns. The function
// returned can be called to stop iteration early. If iteration
// is allowed to complete normally, the function need not be called.
// Either way, the channel will be closed when iteration has finished.
type Iterator func() (in <-chan files.Path, cancel func(), err error)

// callbackIterator is a helper that provides a simple way to
// implement the Iterator interface. The callback will be
// called repeatedly to fetch file paths until either its second
// return parameter is false, it returns an error, or the close function
// is called. The close function is appropriate as a return value for
// the Iterate method itself.
type callbackIterator struct {
	callback func() (files.Path, bool, error)
	cutoff   chan interface{}
}

func (i *callbackIterator) iterate() (<-chan files.Path, func(), error) {
	i.cutoff = make(chan interface{})
	dest := make(chan files.Path)

	// Read files into a channel so that we can select over the
	// inputs and the retries
	go func() {
		for {
			select {
			case _, more := <-i.cutoff:
				if !more {
					close(dest)
					i.cutoff = nil
					return
				}
			default:
				inPath, more, err := i.callback()
				if err != nil {
					slog.Error("iterator callback error: %v", err)
					i.close()
					continue
				}
				if !more {
					i.close()
					continue
				}
				dest <- inPath
			}
		}
	}()

	return dest, i.close, nil
}

func (i *callbackIterator) close() {
	close(i.cutoff)
}
