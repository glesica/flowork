package inputs

import (
	"github.com/glesica/flowork/internal/pkg/files"
	"path/filepath"
	"regexp"
	"strings"
)

type Filter func(p files.Path) bool

func WithRegexp(e string) Filter {
	r := regexp.MustCompile(e)
	return func(p files.Path) bool {
		return r.MatchString(string(p))
	}
}

func WithExt(e string) Filter {
	if !strings.HasPrefix(e, ".") {
		e = "." + e
	}
	return func(p files.Path) bool {
		ae := filepath.Ext(string(p))
		return e == ae
	}
}
