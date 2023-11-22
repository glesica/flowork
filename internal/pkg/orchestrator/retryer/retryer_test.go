package retryer

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/glesica/flowork/internal/pkg/orchestrator"
)

// Implement our own timeout since bugs with channels
// often end up causing the code to hang.
var timeout = 3 * time.Second

func TestNewWorker(t *testing.T) {
	t.Run("should retry a job with an error", func(t *testing.T) {
		errorQueue := make(chan *orchestrator.Job)
		retryQueue := make(chan *orchestrator.Job)

		err := Start(errorQueue, retryQueue, WithMaxRetries(1), WithMaxFailures(math.MaxInt))
		assert.NoError(t, err)

		errorQueue <- &orchestrator.Job{Attempts: 1, Err: errors.New("error")}

		select {
		case retryJob, more := <-retryQueue:
			assert.True(t, more)
			assert.NotZero(t, retryJob)
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	})

	t.Run("should not retry a job that has exceeded max retries", func(t *testing.T) {
		errorQueue := make(chan *orchestrator.Job)
		retryQueue := make(chan *orchestrator.Job)

		err := Start(errorQueue, retryQueue, WithMaxRetries(1), WithMaxFailures(math.MaxInt))
		assert.NoError(t, err)

		errorQueue <- &orchestrator.Job{Attempts: 2, Err: errors.New("error")}
		close(errorQueue)

		select {
		case retryJob, more := <-retryQueue:
			assert.False(t, more)
			assert.Zero(t, retryJob)
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	})

	t.Run("should not retry a job if max failures has been exceeded", func(t *testing.T) {
		errorQueue := make(chan *orchestrator.Job)
		retryQueue := make(chan *orchestrator.Job)

		err := Start(errorQueue, retryQueue, WithMaxRetries(1), WithMaxFailures(0))
		assert.NoError(t, err)

		// This one has now failed
		errorQueue <- &orchestrator.Job{Attempts: 2, Err: errors.New("error")}
		// This one won't be retried
		errorQueue <- &orchestrator.Job{Attempts: 1, Err: errors.New("error")}
		close(errorQueue)

		select {
		case retryJob, more := <-retryQueue:
			assert.False(t, more)
			assert.Zero(t, retryJob)
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	})
}
