package task

import (
	"log/slog"
)

func DockerRun(inst *Instance, v Volume, user string) ([]string, error) {
	containerWorkDir := inst.GetWorkDir()

	command := []string{
		"docker",
		"run",
		"--rm",
		// "--read-only",
	}

	// Set necessary volumes (-v)
	command = append(command, "-v", string(v)+":"+containerWorkDir)

	// Set working directory (-w)
	command = append(command, "-w", containerWorkDir)

	// Set user:group (-u)
	command = append(command, "-u", user+":"+user)

	// Set environment variables (-e)
	// TODO: Implement environment variables

	// Set image
	command = append(command, inst.Image)

	// TODO: Support multiple commands
	command = append(command, inst.Cmd...)

	slog.Debug("running command", "command", command)

	return command, nil
}
