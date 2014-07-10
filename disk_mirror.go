package tvdb

import "path"
import "io"
import "os"
import "log"
import "reflect"

type DiskMirror struct {
	Path           string
	DirPermissions os.FileMode
	Permissions    os.FileMode
}

func NewDiskMirror(path string, dirPermissions os.FileMode, permissions os.FileMode) (diskMirror DiskMirror, err error) {
	err = os.MkdirAll(path, dirPermissions)
	if err != nil {
		return
	}
	diskMirror = DiskMirror{Path: path, DirPermissions: dirPermissions, Permissions: permissions}
	return
}

func (diskMirror DiskMirror) Get(relativePath string, api ApiClient) (reader io.ReadCloser, err error) {
	absPath := path.Join(diskMirror.Path, relativePath)
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		log.Println("DiskMirror: caching " + absPath)
		// If we don't already have it, get it from the server
		var apiReader io.ReadCloser
		apiReader, err = api.Get(relativePath)
		if err != nil {
			return
		}
		defer apiReader.Close()

		// Then create the directory which will house it (if necessary)
		_, err = os.Stat(path.Dir(absPath))
		if os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(absPath), diskMirror.DirPermissions)
		}
		if err != nil {
			return
		}

		// Then create the file
		var file *os.File
		file, err = os.Create(absPath)
		if err != nil {
			return
		}
		defer file.Close()

		err = file.Chmod(diskMirror.Permissions)
		if err != nil {
			return
		}

		// Write to the file
		_, err = io.Copy(file, apiReader)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	} else {
		log.Println("DiskMirror: reading cached file " + absPath)
	}

	// Finally, read the file
	return os.Open(absPath)
}

func (diskMirror DiskMirror) GetXml(relativePath string, api ApiClient, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = diskMirror.Get(relativePath, api)
	if err != nil {
		return
	}
	defer reader.Close()
	return xmlDecode(reader, result)
}

func (diskMirror DiskMirror) GetXmlList(relativePath string, api ApiClient, elementName string, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = diskMirror.Get(relativePath, api)
	if err != nil {
		return
	}
	defer reader.Close()
	err = xmlDecodeList(reader, elementName, result)
	absPath := path.Join(diskMirror.Path, relativePath)
	log.Printf("DiskMirror: %s has %d results", absPath, reflect.ValueOf(result).Elem().Len())
	return
}
