package inputs

import (
	"testing"
)

func TestLocal(t *testing.T) {
	t.Run("unfiltered", func(t *testing.T) {
		// iter, err := Local(files.Dir("fixtures"))
		// assert.NoError(t, err)
		//
		// for _, tc := range []struct {
		// 	value files.Path
		// 	ok    bool
		// }{
		// 	{"fixtures/file0.txt", true},
		// 	{"fixtures/file1.txt", true},
		// 	{"fixtures/file2.txt", true},
		// 	{"fixtures/file3.txt", true},
		// 	{"", false},
		// } {
		// 	//
		// }
	})

	t.Run("filtered", func(t *testing.T) {
		// iter, err := Local("fixtures", func(p files.Path) bool {
		// 	return strings.Contains(string(p), "1")
		// })
		// assert.NoError(t, err)
		//
		// for _, tc := range []struct {
		// 	value files.Path
		// 	ok    bool
		// }{
		// 	{"fixtures/file1.txt", true},
		// 	{"", false},
		// } {
		// 	value, ok := iter()
		// 	assert.Equal(t, tc.value, value)
		// 	assert.Equal(t, tc.ok, ok)
		// }
	})
}
