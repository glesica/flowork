package task

import (
	"github.com/glesica/flowork/internal/pkg/files"
)

// Volume is an opaque string that allows a runner implementation to
// refer to a data volume of some sort, whether that is a local system
// path, a mountable volume on a cloud infrastructure provider, or
// something else. A volume should never include a file name, it must
// be a generic location that can be mounted.
type Volume string

// Runner is the runner interface that allows different
// backends to execute workflow tasks.
type Runner interface {
	// CreateVolume creates a working volume to be used with one or more
	// task instances. The volume must be at least as large as the
	// given Size.
	//
	// Data will be copied to the volume automatically.
	CreateVolume(s files.Size) (Volume, error)

	// DeleteVolume deletes the given volume and all of its contents.
	// It will be called at some point after a task instance has
	// completed and its outputs have been recovered, likely before the
	// full workflow has finished.
	DeleteVolume(v Volume) error

	// AddFile copies the file stored at the given path to the given
	// volume, by whatever means makes the most sense. The file should
	// be copied to the root of the volume, with the given name. The
	// volume reference should not include a file name.
	AddFile(s files.Path, v Volume, name string) error

	// ExtractFile copies a source file from a volume to a different,
	// external path. It is used to recover error logs and outputs
	// from tasks. The destination should not include a file name,
	// the name will be taken from the source path.
	ExtractFile(s files.Path, v Volume, d files.Dir) error

	// Run executes the given task using the given data volume as the
	// task working directory.
	Run(t Instance, v Volume) error
}
