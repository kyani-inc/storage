package bolt

import (
	"fmt"
	"testing"
)

func TestBolt(t *testing.T) {
	a := "Hello Bolt Storage"

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

	f.Flush()
}

func BenchmarkBoltWrite(b *testing.B) {
	f := New("/tmp/benchwrite")

	for n := 0; n < b.N; n++ {
		str := fmt.Sprint("%d", n)

		err := f.Put(str, []byte(str))
		if err != nil {
			b.Errorf("error on put: %s\n", err.Error())
		}
	}

	f.Flush()
}

func BenchmarkBoltRead(b *testing.B) {
	f := New("/tmp/benchread")

	f.Put(string(1), []byte("HELLO WORLD"))

	for n := 0; n < b.N; n++ {
		f.Get(string(1))
	}

	f.Flush()
}
