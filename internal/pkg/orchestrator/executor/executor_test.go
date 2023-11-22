package executor

import (
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/orchestrator"
)

func TestStart(t *testing.T) {
	t.Run("should send a successful job to the success queue", func(t *testing.T) {
		input := make(chan files.Path, 1)
		sucQueue := make(chan *orchestrator.Job, 1)
		defer close(input)

		err := Start(input, WithSuccessQueue(sucQueue))
		assert.NoError(t, err)

		input <- "foo"

		job, more := <-sucQueue
		assert.True(t, more)
		assert.NotZero(t, job)
		assert.Zero(t, job.Err)
		assert.Equal(t, "foo", job.InPath)
	})

	t.Run("should send an errored job to the error queue", func(t *testing.T) {
		input := make(chan files.Path, 1)
		errQueue := make(chan *orchestrator.Job, 1)
		defer close(input)

		err := Start(input, WithErrorQueue(errQueue), WithEngine(func(job *orchestrator.Job) error {
			return errors.New("error")
		}))
		assert.NoError(t, err)

		input <- "foo"

		job, more := <-errQueue
		assert.True(t, more)
		assert.NotZero(t, job)
		assert.NotZero(t, job.Err)
		assert.Equal(t, "foo", job.InPath)
	})

	t.Run("should close success queue on shutdown", func(t *testing.T) {

	})

	t.Run("should close error queue on shutdown", func(t *testing.T) {

	})
}
