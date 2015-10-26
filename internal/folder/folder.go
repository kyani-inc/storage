package folder

import (
	"io/ioutil"
	"strings"
	"sync"
	"os"
	"path/filepath"
)

var lock sync.Mutex

type Folder struct {
	Path string
}

func New(path string) (Folder, error) {
	_, err := os.Stat(path)

	return Folder{Path: path}, err
}

func (f Folder) Put(fileName string, data []byte) error {
	lock.Lock()
	defer lock.Unlock()
	return ioutil.WriteFile(location(f.Path, fileName), data, 0644)
}

func (f Folder) Get(fileName string) []byte {
	data, _ := ioutil.ReadFile(location(f.Path, fileName))
	return data
}

func (f Folder) Delete(fileName string) {
	// ToDo: implement
}

func (f Folder) Flush() {
	filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
		return os.Remove(info.Name())
	})
}

func location(path, fileName string) string {
	return strings.TrimRight(path, "/") + "/" + strings.TrimLeft(fileName, "/")
}
