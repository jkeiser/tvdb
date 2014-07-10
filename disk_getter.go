package tvdb

import "io"
import "os"
import "path"

type DiskGetter struct {
	Root string
}

func (getter DiskGetter) Path(relativePath string) string {
	return path.Join(getter.Root, relativePath)
}

func (getter DiskGetter) Get(relativePath string) (reader io.ReadCloser, err error) {
	return os.Open(getter.Path(relativePath))
}
