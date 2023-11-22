package files

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// Local provides a Store interface into the local file system. All
// paths passed to its methods must be absolute (must begin with /)
// or they will not be accepted.
type Local struct{}

func (l *Local) Accepts(p Path) bool {
	return l.accepts(p) == nil
}

func (l *Local) accepts(p Path) error {
	if !strings.HasPrefix(string(p), "/") {
		return fmt.Errorf("path %s is not absolute", p)
	}

	return nil
}

func (l *Local) Load(p Path) (io.ReadCloser, error) {
	if err := l.accepts(p); err != nil {
		return nil, fmt.Errorf("Local.Load: %w", err)
	}
	return os.Open(string(p))
}

func (l *Local) Save(p Path, f io.Reader) error {
	if err := l.accepts(p); err != nil {
		return fmt.Errorf("Local.Save: %w", err)
	}

	fileDir := path.Dir(string(p))
	if fileDir != "" {
		err := os.MkdirAll(fileDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory (%s): %w", fileDir, err)
		}
	}

	file, err := os.Create(string(p))
	if err != nil {
		return fmt.Errorf("failed to save local file %s: %w", p, err)
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, f)
	if err != nil {
		return fmt.Errorf("failed to copy local file %s: %w", p, err)
	}

	return nil
}

func (l *Local) Close() error {
	return nil
}
