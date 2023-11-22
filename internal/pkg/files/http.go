package files

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/glesica/flowork/internal/pkg/option"
)

type Http struct {
	client *http.Client
}

func NewHttp(opts ...option.Func[*Http]) (*Http, error) {
	s := &Http{}
	err := option.Apply(s, opts...)
	if err != nil {
		return nil, fmt.Errorf("NewHttp: failed to apply options: %w", err)
	}

	if s.client == nil {
		s.client = http.DefaultClient
	}

	return s, nil
}

func WithClient(client *http.Client) option.Func[*Http] {
	return func(s *Http) error {
		s.client = client
		return nil
	}
}

func (h *Http) Accepts(p Path) bool {
	u := string(p)
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

func (h *Http) Load(p Path) (io.ReadCloser, error) {
	resp, err := h.client.Get(string(p))
	if err != nil {
		return nil, fmt.Errorf("Http.Load: error fetching %s: %w", p, err)
	}

	return resp.Body, nil
}

func (h *Http) Save(p Path, f io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (h *Http) Close() error {
	h.client.CloseIdleConnections()
	return nil
}
