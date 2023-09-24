package files

import (
	"path/filepath"
)

// Path is a reference to a file that can exist in any
// supported storage environment. It is abstracted here for
// use with the Store and Runner interfaces. In the future, the
// implementation may need to change to support paths that
// refer to different environments (such as S3).
type Path string

// Dir returns the directory portion of the given path. Like
// Path, it is abstracted for use with the Store and Runner
// interfaces.
func (p Path) Dir() Dir {
	return Dir(filepath.Dir(string(p)))
}

// File returns the file name portion of the given path.
func (p Path) File() string {
	return filepath.Base(string(p))
}

// Dir is a reference to a directory (or similar concept) that
// can exist in any supported storage environment.
type Dir string

// PathTo returns a path to the named file in this directory.
func (d Dir) PathTo(name string) Path {
	return Path(filepath.Join(string(d), name))
}

func (d Dir) SubDir(name string) Dir {
	return Dir(filepath.Join(string(d), name))
}
