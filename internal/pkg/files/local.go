package files

import (
	"fmt"
	"io"
	"os"
	"path"
)

type LocalStore struct{}

func (l *LocalStore) Accepts(s Path) bool {
	return true
}

func (l *LocalStore) Load(p Path) (io.ReadCloser, error) {
	return os.Open(string(p))
}

func (l *LocalStore) Save(p Path, reader io.Reader) error {
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

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to copy local file %s: %w", p, err)
	}

	return nil
}
