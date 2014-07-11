package tvdb

import "path"
import "io"
import "os"
import "log"

type DiskCache struct {
	Root             string
	DefaultSource    RelativePathGetter
	ExtensionSources map[string]RelativePathGetter
	DirPermissions   os.FileMode
	Permissions      os.FileMode
}

func (diskCache DiskCache) sourceFor(relativePath string) RelativePathGetter {
	ext := path.Ext(relativePath)
	getter, ok := diskCache.ExtensionSources[ext]
	if !ok {
		getter = diskCache.DefaultSource
	}
	return getter
}

func (diskCache DiskCache) Path(relativePath string) string {
	return path.Join(diskCache.Root, relativePath)
}

func (diskCache DiskCache) Get(relativePath string) (reader io.ReadCloser, err error) {
	absPath := diskCache.Path(relativePath)
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		diskCache.update(relativePath, absPath)
	} else if err != nil {
		return
	} else {
		log.Println("DiskRepository: reading cached file " + absPath)
	}

	// Finally, read the file
	return os.Open(absPath)
}

func (diskCache DiskCache) update(relativePath string, absPath string) (err error) {
	log.Println("DiskCache: caching " + absPath)
	// If we don't already have it, get it from the server
	var reader io.ReadCloser
	reader, err = diskCache.sourceFor(relativePath).Get(relativePath)
	if err != nil {
		return
	}
	defer reader.Close()

	// Then create the directory which will house it (if necessary)
	_, err = os.Stat(path.Dir(absPath))
	if os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(absPath), diskCache.DirPermissions)
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

	err = file.Chmod(diskCache.Permissions)
	if err != nil {
		return
	}

	// Write to the file
	_, err = io.Copy(file, reader)
	return
}
