package memcache_test
import (
"testing"
"github.com/kyani-inc/storage/internal/memcache"
	"strings"
	"os"
)

func getHosts() []string {
	v := os.Getenv("MEMCACHE_HOSTS")

	if v == "" {
		return []string{}
	}

	return strings.Split(v, ",")
}

func TestMemcache(t *testing.T) {
	hosts := getHosts()

	if len(hosts) < 1 {
		t.Skip("need memcache hosts to test")
	}

	k, v := "greeting", "Hello World"

	m := memcache.New(hosts)

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
