package files

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"cloud.google.com/go/storage"

	"github.com/glesica/flowork/internal/pkg/option"
)

// Gcs provides a Store interface that can handle files stored in
// Google Cloud Storage (GCS).
type Gcs struct {
	client *storage.Client
}

func NewGcs(opts ...option.Func[*Gcs]) (*Gcs, error) {
	s := &Gcs{}
	err := option.Apply(s, opts...)
	if err != nil {
		return nil, fmt.Errorf("NewGcs: failed to apply options: %w", err)
	}

	// TODO: Create a client if one wasn't provided

	return s, nil
}

// WithGcsClient provides a GCS client to use. The Gcs will
// assume ownership of this client, so it should not be modified
// or closed from the outside. Calling Close will also close the
// client.
func WithGcsClient(client *storage.Client) option.Func[*Gcs] {
	return func(s *Gcs) error {
		s.client = client
		return nil
	}
}

func (s *Gcs) Accepts(p Path) bool {
	u, err := url.Parse(string(p))
	if err != nil {
		return false
	}

	if u.Scheme != "gs" {
		return false
	}

	if u.Path == "" {
		return false
	}

	return true
}

func (s *Gcs) Load(p Path) (io.ReadCloser, error) {
	u, err := url.Parse(string(p))
	if err != nil {
		return nil, fmt.Errorf("Gcs.Load: failed to parse gs url: %w", err)
	}

	b := s.client.Bucket(u.Host)
	o := b.Object(u.Path)

	return o.NewReader(context.Background())
}

func (s *Gcs) Save(p Path, f io.Reader) error {
	u, err := url.Parse(string(p))
	if err != nil {
		return fmt.Errorf("Gcs.Save: failed to parse gs url: %w", err)
	}

	b := s.client.Bucket(u.Host)
	o := b.Object(u.Path)
	r := o.NewWriter(context.Background())
	defer func() { _ = r.Close() }()

	_, err = io.Copy(r, f)
	if err != nil {
		return fmt.Errorf("Gcs.Save: failed to copy file: %w", err)
	}

	return nil
}

func (s *Gcs) Close() error {
	err := s.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Gcs: %w", err)
	}

	return nil
}
