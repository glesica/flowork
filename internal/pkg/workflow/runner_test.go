package workflow

import (
	"github.com/alecthomas/assert/v2"
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/inputs"
	"github.com/glesica/flowork/internal/pkg/spec"
	"github.com/glesica/flowork/internal/pkg/task"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Success(t *testing.T) {
	t.Cleanup(func() {
		// Comment this out to preserve the outputs after the test run
		_ = os.RemoveAll("fixtures/output")
	})

	ws, err := spec.LoadWorkflowPath("fixtures/workflow.json")
	assert.NoError(t, err)
	wi := NewInstance(ws)

	store := &files.LocalStore{}

	workDir, err := filepath.Abs("fixtures/")
	assert.NoError(t, err)

	runner := &task.DockerRunner{
		Debug:   false,
		WorkDir: files.Dir(workDir),
		Store:   store,
	}

	inFiles, err := inputs.Local("fixtures/inputs")
	assert.NoError(t, err)

	err = Run(wi, runner, inFiles, "fixtures/output", 1)
	assert.NoError(t, err)

	fi, err := os.Stat("fixtures/outputs")
	assert.NoError(t, err)
	assert.True(t, fi.IsDir())
}
