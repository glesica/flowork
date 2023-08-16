package task

import (
	"fmt"
	"github.com/glesica/flowork/internal/pkg/shell"
	"io"
	"os"
	"os/user"
	"path"
)

type DockerRunner struct{}

func (d *DockerRunner) Run(inst Instance) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	hostWorkDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	containerWorkDir := inst.GetWorkDir()

	command := []string{
		"docker",
		"run",
		"--rm",
		"--read-only",
	}

	// Set necessary volumes (-v)
	command = append(command, "-v", hostWorkDir+":"+containerWorkDir)

	// Set working directory (-w)
	command = append(command, "-w", containerWorkDir)

	// Set user:group (-u)
	command = append(command, "-u", currentUser.Uid+":"+currentUser.Gid)

	// Set environment variables (-e)
	// TODO: Implement environment variables

	// Set image
	command = append(command, inst.Image)

	// TODO: Support multiple commands
	command = append(command, inst.Cmd...)

	result, err := shell.Run(command)
	if err != nil {
		return fmt.Errorf("failed to run docker (%v): %w", command, err)
	}

	err = writeOutput(hostWorkDir, "stdout.txt", result.Out)
	if err != nil {
		return fmt.Errorf("failed to write stdout.txt: %w", err)
	}

	err = writeOutput(hostWorkDir, "stderr.txt", result.Err)
	if err != nil {
		return fmt.Errorf("failed to write stderr.txt: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("task failed, see stderr.txt")
	}

	return nil
}

func writeOutput(workDir, name, data string) error {
	outPath := path.Join(workDir, name)
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
