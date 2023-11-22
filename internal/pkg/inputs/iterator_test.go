package inputs

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/glesica/flowork/internal/pkg/files"
)

var paths = []files.Path{"foo", "bar", "baz"}

func getCallback() func() (files.Path, bool, error) {
	next := 0
	return func() (files.Path, bool, error) {
		if next >= len(paths) {
			return "", false, nil
		}
		p := paths[next]
		next++
		return p, true, nil
	}
}

func Test_callbackIterator_iterate(t *testing.T) {
	cbi := callbackIterator{callback: getCallback()}
	pc, _, _ := cbi.iterate()

	assert.Equal(t, "foo", <-pc)
	assert.Equal(t, "bar", <-pc)
	assert.Equal(t, "baz", <-pc)

	p, more := <-pc
	assert.Equal(t, "", p)
	assert.False(t, more)
}

func Test_callbackIterator_close(t *testing.T) {
	cbi := callbackIterator{callback: getCallback()}
	pc, _, _ := cbi.iterate()

	assert.Equal(t, "foo", <-pc)
	cbi.close()
	assert.Equal(t, "bar", <-pc)

	p, more := <-pc
	assert.Equal(t, "", p)
	assert.False(t, more)
}
