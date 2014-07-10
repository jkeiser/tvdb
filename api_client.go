package tvdb

import "net/http"
import "io"
import "log"
import "reflect"

type ApiClient struct {
	Url string
}

/*
 * Helpers
 */
func (api ApiClient) url(relative string) (url string) {
	return api.Url + "/" + relative
}

func (api ApiClient) Get(relative string) (reader io.ReadCloser, err error) {
	log.Println("ApiClient: GET " + api.url(relative))
	response, err := http.Get(api.url(relative))
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	reader = response.Body
	return reader, nil
}

func (api ApiClient) GetXml(relative string, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = api.Get(relative)
	if err != nil {
		return
	}
	defer reader.Close()
	return xmlDecode(reader, result)
}

func (api ApiClient) GetXmlList(relative string, elementName string, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = api.Get(relative)
	if err != nil {
		return
	}
	defer reader.Close()
	err = xmlDecodeList(reader, elementName, result)
	log.Printf("ApiClient: GET %s returned %d results", api.url(relative), reflect.ValueOf(result).Elem().Len())
	return
}
