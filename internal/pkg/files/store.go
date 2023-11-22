package files

import (
	"io"
)

// TODO: Add a way to find file size
// TODO: Should paths be URIs?

type Store interface {
	// Accepts indicates whether a given store can operate on the
	// given path. It might do this, for example, by checking its
	// prefix for a protocol. If the Store claims to accept a Path,
	// then that Path will be assumed to point at a file (whether
	// the file actually exists is a separate matter).
	Accepts(p Path) bool

	// Load opens the file represented by the given path, whatever
	// that means in the context of a given store, and returns a
	// suitable reader.
	Load(p Path) (io.ReadCloser, error)

	// Save writes the data taken from the given reader to a file
	// stored at the given path.
	Save(p Path, f io.Reader) error

	// Close renders the store unusable by closing or otherwise
	// disposing any resources it was using. It may block until this
	// process is complete.
	Close() error
}
