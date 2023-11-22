package inputs

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/glesica/flowork/internal/pkg/files"
)

func TestWithRegexp(t *testing.T) {
	for _, c := range []struct {
		name    files.Path
		regex   string
		outcome bool
	}{
		{"/dir/file", `^/dir/[a-z]+$`, true},
	} {
		f := WithRegexp(c.regex)
		assert.Equal(t, c.outcome, f(c.name))
	}
}

func TestWithExt(t *testing.T) {
	for _, c := range []struct {
		name    files.Path
		ext     string
		outcome bool
	}{
		{"/dir/file.foo", "foo", true},
		{"/dir/file.foo.bar", "foo", false},
		{"/dir/filefoo", "foo", false},
		{"foo", "foo", false},
		{"file.foo", "foo", true},
	} {
		f := WithExt(c.ext)
		assert.Equal(t, c.outcome, f(c.name))
	}
}
