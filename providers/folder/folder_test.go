package folder_test

import (
	"testing"

	"github.com/kyani-inc/storage/providers/folder"
	"io/ioutil"
)

func TestFolder(t *testing.T) {
	p, err := ioutil.TempDir("", "storage-folder-test")

	if err != nil {
		t.Fatal(err.Error())
	}

	k, v := "folder_test", "Hello Folder Storage"

	f, err := folder.New(p)

	if err != nil {
		t.Fatal(err.Error())
	}

	err = f.Put(k, []byte(v))

	if err != nil {
		t.Fatalf("error on put: %s", err.Error())
	}

	b := f.Get(k)

	if v != string(b) {
		t.Fatalf("expected '%s' but got '%s'", v, b)
	}

	f.Delete(k)

	c := f.Get(k)

	if len(c) > 0 {
		t.Errorf("expected emtpy result; got %s", c)
	}

	f.Flush()
}
