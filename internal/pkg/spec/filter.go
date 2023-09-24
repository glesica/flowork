package spec

import (
	"fmt"
	"regexp"
)

type Filter struct {
	// Path provides a regular expression, matched against the full,
	// absolute file path. Only paths that match will be passed through.
	//
	// Ex:
	//   - /home/bob/data/file.csv
	//   - aws://bobs-bucket/data/file.csv
	//   - https://bob.net/data/file.csv
	Path string `json:"path"`
}

type FilterError struct {
	Message string
	Filter  *Filter
	URI     string
	Wrapped error
}

func (e *FilterError) Error() string {
	return fmt.Sprintf("%s: filter: %+v uri: %s", e.Message, e.Filter, e.URI)
}

func (e *FilterError) Unwrap() error {
	return e.Wrapped
}

func (f *Filter) Match(uri string) (bool, error) {
	pathMatch, err := regexp.MatchString(f.Path, uri)
	if err != nil {
		return false, &FilterError{
			Message: "path regex failed",
			Filter:  f,
			URI:     uri,
			Wrapped: err,
		}
	}
	if pathMatch {
		return true, nil
	}

	return false, nil
}
