package inputs

import (
	"fmt"
	"os"

	"github.com/glesica/flowork/internal/pkg/files"
)

// Local provides an iterator over all the normal files in a given
// directory. It does not traverse into subdirectories. The given
// path must represent a directory.
//
// TODO: This should be called Dir or something since it lists files in a directory
// It should also work with any of the stores that support listing files (or will?)
func Local(dir files.Dir, filters ...Filter) (Iterator, error) {
	entries, err := os.ReadDir(string(dir))
	if err != nil {
		return nil, fmt.Errorf("failed to get local inputs (%s): %w", dir, err)
	}

	index := 0

	cbi := &callbackIterator{
		callback: func() (files.Path, bool, error) {
			for {
				if index >= len(entries) {
					return "", false, nil
				}

				e := entries[index]
				index++

				name := e.Name()
				ep := dir.PathTo(name)

				accept := true
				for _, f := range filters {
					if !f(ep) {
						accept = false
						break
					}
				}

				if accept {
					return dir.PathTo(name), true, nil
				}
			}
		},
		cutoff: make(chan interface{}),
	}

	return cbi.iterate, nil
}
