package folder

import (
	"io/ioutil"
	"strings"
	"sync"
)

var lock sync.Mutex

type Folder struct {
	Path string
}

func New(path string) Folder {
	return Folder{Path: path}
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

func location(path, fileName string) string {
	return strings.TrimRight(path, "/") + "/" + strings.TrimLeft(fileName, "/")
}
