package files

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestPath_Dir(t *testing.T) {
	t.Run("should extract dir", func(t *testing.T) {
		p := Path("/a/b/c/d.txt")
		d := p.Dir()
		assert.Equal(t, "/a/b/c", d)
	})

	t.Run("should return an empty dir", func(t *testing.T) {
		p := Path("d.txt")
		d := p.Dir()
		assert.Equal(t, ".", d)
	})
}

func TestPath_File(t *testing.T) {
	t.Run("should extract file name", func(t *testing.T) {
		p := Path("/a/b/c/d.txt")
		f := p.File()
		assert.Equal(t, "d.txt", f)
	})

	t.Run("should return a bare file name", func(t *testing.T) {
		p := Path("d.txt")
		f := p.File()
		assert.Equal(t, "d.txt", f)
	})
}

func TestDir_PathTo(t *testing.T) {
	t.Run("should append file name", func(t *testing.T) {
		d := Dir("/a/b/c")
		p := d.PathTo("d.txt")
		assert.Equal(t, "/a/b/c/d.txt", p)
	})
}

func TestDir_SubDir(t *testing.T) {
	t.Run("should append dir name", func(t *testing.T) {
		d := Dir("/a/b/c")
		s := d.SubDir("d")
		assert.Equal(t, "/a/b/c/d", s)
	})
}
