package workflow

import (
	"errors"
	"fmt"
	"sync"
)

type ErrSet struct {
	mutex *sync.Mutex
	errs  []error
}

func NewErrSet() *ErrSet {
	return &ErrSet{
		mutex: &sync.Mutex{},
		errs:  nil,
	}
}

func (e *ErrSet) AddFormat(format string, a ...any) {
	err := fmt.Errorf(format, a...)
	e.Add(err)
}

func (e *ErrSet) Add(err error) {
	e.mutex.Lock()
	e.errs = append(e.errs, err)
	e.mutex.Unlock()
}

func (e *ErrSet) Count() int {
	e.mutex.Lock()
	c := len(e.errs)
	e.mutex.Unlock()

	return c
}

func (e *ErrSet) Capture(f func() error) {
	go func() {
		err := f()
		if err != nil {
			e.Add(err)
		}
	}()
}

func (e *ErrSet) Err() error {
	e.mutex.Lock()
	err := errors.Join(e.errs...)
	e.mutex.Unlock()

	return err
}
