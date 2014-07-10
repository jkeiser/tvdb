package tvdb

import "io"

type PathGetter interface {
	Get(relativePath string) (reader io.ReadCloser, err error)
	Path(relativePath string) (path string)
}
