package shell

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
)

type Result struct {
	Code int
	Out  string
	Err  string
}

func Run(cmd []string) (*Result, error) {
	slog.Debug("running shell command", "command", cmd)
	c := exec.Command(cmd[0], cmd[1:]...)

	outBuf, errBuf := bytes.Buffer{}, bytes.Buffer{}
	c.Stdout = &outBuf
	c.Stderr = &errBuf

	err := c.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		r := &Result{
			Code: exitErr.ExitCode(),
			Out:  outBuf.String(),
			Err:  errBuf.String(),
		}
		return r, fmt.Errorf("failed to run command (%s): %w", c.String(), exitErr)
	} else if err != nil {
		return nil, err
	}

	r := &Result{
		Code: 0,
		Out:  outBuf.String(),
		Err:  errBuf.String(),
	}
	slog.Debug("shell command complete", "result", r)

	return r, nil
}
