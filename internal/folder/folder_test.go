package folder

import (
	"sync"
	"testing"
)

func TestFolder(t *testing.T) {
	a := "Hello Folder Storage"

	f := New("/tmp/storage")

	err := f.Put("folder_test", []byte(a))

	if err != nil {
		t.Errorf("error on put: %s", err.Error())
		return
	}

	b := f.Get("folder_test")

	if a != string(b) {
		t.Errorf("expected '%s' but got '%s'", a, b)
	}
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
