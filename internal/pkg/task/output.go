package task

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/glesica/flowork/internal/pkg/files"
)

func outputWriters(outputDir files.Dir) (stdout io.WriteCloser, stderr io.WriteCloser, err error) {
	stdoutPath := outputDir.PathTo("stdout.txt")
	stderrPath := outputDir.PathTo("stderr.txt")

	stdout, err = os.Create(string(stdoutPath))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stdout file: %w", err)
	}

	stderr, err = os.Create(string(stderrPath))
	if err != nil {
		_ = stdout.Close()
		return nil, nil, fmt.Errorf("failed to open stderr file: %w", err)
	}

	return stdout, stderr, nil
}

func writeOutput(outputDir, name, data string) error {
	outPath := path.Join(outputDir, name)
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to open output file (%s): %w", outPath, err)
	}

	_, err = io.WriteString(out, data)
	if err != nil {
		return fmt.Errorf("failed to write output file (%s): %w", outPath, err)
	}

	return nil
}
