package folder_test

import (
	"sync"
	"testing"

	"github.com/kyani-inc/storage/internal/folder"
	"io/ioutil"
)

func TestFolder(t *testing.T) {
	p, err := ioutil.TempDir("", "storage-folder-test")

	if err != nil {
		t.Fatal(err.Error())
	}

	a := "Hello Folder Storage"

	f, err := folder.New(p)

	if err != nil {
		t.Fatal(err.Error())
	}

	err = f.Put("folder_test", []byte(a))

	if err != nil {
		t.Fatalf("error on put: %s", err.Error())
	}

	b := f.Get("folder_test")

	if a != string(b) {
		t.Fatalf("expected '%s' but got '%s'", a, b)
	}

	f.Flush()
}

func TestLocalRace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		TestFolder(t)
	}(&wg, t)

	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		TestFolder(t)
	}(&wg, t)

	wg.Wait()
}
