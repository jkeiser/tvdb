package tvdb

import "path"
import "io"

type RelativeGetter struct {
	RelativeTo   PathGetter
	RelativePath string
}

func (getter RelativeGetter) Get(relativePath string) (reader io.ReadCloser, err error) {
	newPath := path.Join(getter.RelativePath, relativePath)
	return getter.Get(newPath)
}

func (getter RelativeGetter) Path(relativePath string) string {
	newPath := path.Join(getter.RelativePath, relativePath)
	return getter.Path(newPath)
}
