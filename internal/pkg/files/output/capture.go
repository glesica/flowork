package output

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/glesica/flowork/internal/pkg/files"
)

// TODO: Need to stream data into the writers from outside, so need a "done" signal
// Empty buffers return EOF on read, so that won't work because the buffer could
// be temporarily empty

func Writers(s files.Store, dest files.Dir) (stdout io.Writer, stderr io.Writer, err error) {
	stdoutPath := dest.PathTo("stdout.txt")
	stderrPath := dest.PathTo("stderr.txt")

	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}

	go func() {
		err := s.Save(stdoutPath, stdoutBuf)
		if err != nil {
			slog.Error("failed to save stdout (%s): %w", stdoutPath, err)
		}
	}()
	go func() {
		err := s.Save(stderrPath, stderrBuf)
		if err != nil {
			slog.Error("failed to save stderr (%s): %w", stderrPath, err)
		}
	}()

	return stdoutBuf, stderrBuf, nil
}
