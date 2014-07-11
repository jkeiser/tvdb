package tvdb

import "net/http"
import "log"
import "io"
import "os"
import "path"

type RelativePathGetter interface {
	Get(relativePath string) (reader io.ReadCloser, err error)
	Path(relativePath string) (path string)
}

type DiskGetter struct {
	Root string
}

func (getter DiskGetter) Path(relativePath string) string {
	return path.Join(getter.Root, relativePath)
}

func (getter DiskGetter) Get(relativePath string) (reader io.ReadCloser, err error) {
	return os.Open(getter.Path(relativePath))
}

type HttpGetter struct {
	Url string
}

func (api HttpGetter) Path(relative string) (url string) {
	return api.Url + "/" + relative
}

func (api HttpGetter) Get(relative string) (reader io.ReadCloser, err error) {
	log.Println("HttpGetter: GET " + api.Path(relative))
	response, err := http.Get(api.Path(relative))
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	reader = response.Body
	return reader, nil
}

type RelativeGetter struct {
	RelativeTo   RelativePathGetter
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
