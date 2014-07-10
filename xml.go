package tvdb

import "encoding/xml"
import "reflect"
import "io"
import "errors"
import "strings"

func GetXml(reader io.Reader, result interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	return decoder.Decode(result)
}

func GetXmlList(reader io.Reader, elementName string, result interface{}) (err error) {
	// Get []<type> (dereference *[]<type>)
	slice := reflect.ValueOf(result).Elem()
	elementType := slice.Type().Elem()

	// Skip to the first start tag (the actual element)
	decoder := xml.NewDecoder(reader)
	var start xml.StartElement
	start, err = xmlSkipUntilStart(decoder)
	if err != nil {
		return
	}

	// Read subsequent elements, constructing and decoding an ElementType for each
	for {
		start, err = xmlSkipUntilStart(decoder)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return
		}
		if start.Name.Local != elementName {
			return errors.New("Found XML element " + start.Name.Local + ", expected only " + elementName)
		}
		// New returns a pointer to the value, so we need to get the value itself
		newValuePtr := reflect.New(elementType)
		decoder.DecodeElement(newValuePtr.Interface(), &start)
		slice.Set(reflect.Append(slice, newValuePtr.Elem()))
	}
}

func xmlSkipUntilStart(decoder *xml.Decoder) (start xml.StartElement, err error) {
	for {
		// Read tokens from the XML document in a stream.
		var token xml.Token
		token, err = decoder.Token()
		if err != nil {
			return
		}

		var ok bool
		start, ok = token.(xml.StartElement)
		if ok {
			return
		}
	}
}

type PipeDelimitedString []string

func (array *PipeDelimitedString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}
	*array = strings.Split(content, "|")
	return nil
}
