package inputs

import (
	"fmt"
	"github.com/glesica/flowork/internal/pkg/files"
	"os"
	"path/filepath"
)

// Local provides an iterator over all the normal files in a given
// directory. It does not traverse into subdirectories. The given
// path must represent a directory.
func Local(p files.Dir, filters ...Filter) (files.Iter, error) {
	entries, err := os.ReadDir(string(p))
	if err != nil {
		return nil, fmt.Errorf("failed to get local inputs (%s): %w", p, err)
	}

	i := 0
	return func() (files.Path, bool) {
		for {
			if i >= len(entries) {
				entries = nil
				return "", false
			}
			e := entries[i]
			i++

			// We only operate on regular files, for now,
			// because it is simpler that way.
			if e.Type() != 0 {
				continue
			}

			ep := files.Path(filepath.Join(string(p), e.Name()))

			accept := true
			for _, f := range filters {
				if !f(ep) {
					accept = false
					break
				}
			}

			if accept {
				return ep, true
			}
		}

	}, nil
}
