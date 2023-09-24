package task

import (
	"fmt"
	"github.com/glesica/flowork/internal/app/options"
	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/id"
	"github.com/glesica/flowork/internal/pkg/shell"
	"io"
	"log/slog"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

type DockerRunner struct {
	// Debug indicates whether debug mode is enabled. In debug mode,
	// the runner doesn't delete volumes and may provide additional
	// output.
	Debug bool

	// WorkDir is the local (host) working directory, used to store
	// data for tasks to operate on. Subdirectories will be created
	// and mounted as volumes in the container that is created to
	// run each task.
	WorkDir files.Dir

	// Store is the store to be used for reading and writing files
	// to volumes.
	Store files.Store
}

func (r *DockerRunner) CreateVolume(s files.Size) (Volume, error) {
	volDir := filepath.Join(string(r.WorkDir), options.VolumesDirName, id.New())

	err := os.MkdirAll(volDir, 0777)
	if err != nil {
		return "", fmt.Errorf("failed to create volume (%s): %w", volDir, err)
	}

	return Volume(volDir), nil
}

func (r *DockerRunner) DeleteVolume(v Volume) error {
	if r.Debug {
		slog.Debug("delete volume requested, ignoring", "volume", v)
		return nil
	}

	err := os.RemoveAll(string(v))
	if err != nil {
		return fmt.Errorf("failed to delete volume %s: %w", v, err)
	}

	return nil
}

// TODO: Make name a path and create intermediate directories

func (r *DockerRunner) AddFile(s files.Path, v Volume, name string) error {
	fileData, err := r.Store.Load(s)
	if err != nil {
		return fmt.Errorf("failed to load file %s for add: %w", s, err)
	}
	defer func() { _ = fileData.Close() }()

	dest := filepath.Join(string(v), name)

	err = r.Store.Save(files.Path(dest), fileData)
	if err != nil {
		return fmt.Errorf("failed to save file %s for add: %w", dest, err)
	}

	return nil
}

func (r *DockerRunner) ExtractFile(s files.Path, v Volume, d files.Dir) error {
	name := s.File()
	src := filepath.Join(string(v), name)

	fileData, err := r.Store.Load(files.Path(src))
	if err != nil {
		return fmt.Errorf("failed to load file %s for extract: %w", s, err)
	}
	defer func() { _ = fileData.Close() }()

	dest := d.PathTo(name)

	err = r.Store.Save(dest, fileData)
	if err != nil {
		return fmt.Errorf("failed to save file %s for extract: %w", d, err)
	}

	return nil
}

func (r *DockerRunner) Run(inst Instance, v Volume) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	containerWorkDir := inst.GetWorkDir()

	command := []string{
		"docker",
		"run",
		"--rm",
		//"--read-only",
	}

	// Set necessary volumes (-v)
	command = append(command, "-v", string(v)+":"+containerWorkDir)

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

	slog.Debug("running command", "command", command)

	result, err := shell.Run(command)
	if err != nil {
		if result != nil {
			_ = writeOutput(string(v), "stdout.txt", result.Out)
			_ = writeOutput(string(v), "stderr.txt", result.Err)
		}
		return fmt.Errorf("failed to run docker (%v): %w", command, err)
	}

	// TODO: Write to workflow and task instance specific directories

	err = writeOutput(string(v), "stdout.txt", result.Out)
	if err != nil {
		return fmt.Errorf("failed to write stdout.txt: %w", err)
	}

	err = writeOutput(string(v), "stderr.txt", result.Err)
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
