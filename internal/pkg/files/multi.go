package files

import (
	"errors"
	"fmt"
	"io"

	"github.com/glesica/flowork/internal/pkg/option"
)

type Multi struct {
	// Various store implementations to be queried.
	stores []Store
}

func NewMulti(opts ...option.Func[*Multi]) (*Multi, error) {
	s := &Multi{}

	err := option.Apply(s, opts...)
	if err != nil {
		return nil, fmt.Errorf("NewMulti: failed to apply options: %w", err)
	}

	return s, nil
}

func WithStore(s Store) option.Func[*Multi] {
	return func(m *Multi) error {
		m.stores = append(m.stores, s)
		return nil
	}
}

func (m *Multi) Accepts(p Path) bool {
	for _, c := range m.stores {
		if c.Accepts(p) {
			return true
		}
	}

	return false
}

func (m *Multi) Load(p Path) (io.ReadCloser, error) {
	for _, c := range m.stores {
		if c.Accepts(p) {
			return c.Load(p)
		}
	}

	return nil, fmt.Errorf("cannot load unsupported path %s", p)
}

func (m *Multi) Save(p Path, data io.Reader) error {
	for _, c := range m.stores {
		if c.Accepts(p) {
			return c.Save(p, data)
		}
	}

	return fmt.Errorf("cannot save unsupported path %s", p)
}

func (m *Multi) Close() error {
	var errs []error
	for _, s := range m.stores {
		err := s.Close()
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
