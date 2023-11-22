package executor

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/glesica/flowork/internal/pkg/spec"
)

func TestMakeJobCreator(t *testing.T) {
	t.Run("should populate fields", func(t *testing.T) {
		createJob := MakeJobCreator(spec.TaskSet{spec.Task{
			Name:    "fake",
			Desc:    "",
			Cmd:     nil,
			Image:   "",
			WorkDir: "",
			Inputs:  nil,
			Outputs: nil,
		}})
		job, err := createJob("/path")

		assert.NoError(t, err)
		assert.NotZero(t, job.Id)
		assert.Equal(t, "/path", job.InPath)
		assert.Equal(t, 1, len(job.Tasks))
		assert.Equal(t, "fake", job.Tasks[0].Name)
	})
}
