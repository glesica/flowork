package task

import (
	"fmt"
	"github.com/glesica/flowork/internal/pkg/shell"
	"os"
	"os/user"
)

type DockerRunner struct{}

func (d *DockerRunner) Name() string {
	return "docker"
}

func (d *DockerRunner) Run(inst Instance) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	command := []string{
		"docker",
		"run",
		"--rm",
		"--read-only",
	}

	t := inst.Task

	// Set necessary volumes (-v)
	command = append(command, "-v", workDir+":"+t.WorkDir)

	// Set working directory (-w)
	command = append(command, "-w", t.WorkDir)

	// Set user:group (-u)
	command = append(command, "-u", currentUser.Uid+":"+currentUser.Gid)

	// Set environment variables (-e)
	// TODO: Implement environment variables

	// Set image
	command = append(command, t.Image)

	// TODO: Support multiple commands
	command = append(command, t.Cmd...)

	result, err := shell.Run(command)
	if err != nil {
		return fmt.Errorf("failed to run docker (%v): %w", command, err)
	}

	if result.Code != 0 {
		return fmt.Errorf("task failed (%s)", result.Err)
	}

	return nil
}
