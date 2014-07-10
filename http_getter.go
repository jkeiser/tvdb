package tvdb

import "net/http"
import "io"
import "log"

type HttpGetter struct {
	Url string
}

/*
 * PathGetter
 */
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
