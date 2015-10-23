package local_test
import (
"testing"
"github.com/kyani-inc/storage/internal/local"
)

func TestLocal(t *testing.T) {
	k, v := "greeting", []byte("Hello World")

	l := local.New()

	err := l.Put(k, v)

	if err != nil {
		t.Fatal(err.Error())
	}

	b := l.Get(k)

	if string(v) != string(b) {
		t.Errorf("expected %s; got %s", v, b)
	}
}
