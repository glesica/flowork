package files

import (
	"fmt"
	"io"
)

type Iter func() (Path, bool)

// TODO: Add a way to find file size

type Store interface {
	// Accepts indicates whether a given store can operate on the
	// given path. It might do this, for example, by checking its
	// prefix for a protocol.
	Accepts(Path) bool

	// Load opens the file represented by the given path, whatever
	// that means in the context of a given store, and returns a
	// suitable reader.
	Load(Path) (io.ReadCloser, error)

	// Save writes the data taken from the given reader to a file
	// stored at the given path.
	Save(Path, io.Reader) error
}

type MetaStore struct {
	// Various store implementations to be queried.
	children []Store
}

func (m *MetaStore) Accepts(p Path) bool {
	for _, c := range m.children {
		if c.Accepts(p) {
			return true
		}
	}

	return false
}

func (m *MetaStore) Load(p Path) (io.ReadCloser, error) {
	for _, c := range m.children {
		if c.Accepts(p) {
			return c.Load(p)
		}
	}

	return nil, fmt.Errorf("cannot load unsupported path %s", p)
}

func (m *MetaStore) Store(p Path, data io.Reader) error {
	for _, c := range m.children {
		if c.Accepts(p) {
			return c.Save(p, data)
		}
	}

	return fmt.Errorf("cannot save unsupported path %s", p)
}
