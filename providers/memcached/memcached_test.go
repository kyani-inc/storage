package memcached_test

import (
	"os"
	"strings"
	"testing"

	"github.com/kyani-inc/storage/providers/memcached"
	"github.com/subosito/gotenv"
)

var hosts []string

func init() {
	// attempt to load env vars in memcache pacakge.
	gotenv.Load(".env")

	v := os.Getenv("MEMCACHED_HOSTS")

	if v != "" {
		hosts = strings.Split(v, ",")
		return
	}
}

func TestMemcached(t *testing.T) {
	if len(hosts) < 1 {
		t.Skip("need memcache hosts to test")
	}

	k, v := "greeting", "Hello World"

	m := memcached.New(hosts)

	err := m.Put(k, []byte(v))

	if err != nil {
		t.Fatal(err.Error())
	}

	b := m.Get(k)

	if v != string(b) {
		t.Errorf("expected %s; got %s", v, b)
	}

	m.Flush()
}
