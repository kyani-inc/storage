package redis_test

import (
	"os"
	"testing"

	"github.com/kyani-inc/storage/providers/redis"
	"github.com/subosito/gotenv"
)

var (
	host, port string
)

func init() {
	gotenv.Load(".env")

	host, port = os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
}

func TestRedis(t *testing.T) {
	if host == "" || port == "" {
		t.Skip("need redis host and port to test")
	}

	k, v := "greetings", "Hello World"

	r := redis.New(host, port)

	err := r.Put(k, []byte(v))

	if err != nil {
		t.Fatal(err.Error())
	}

	b := r.Get(k)

	if v != string(b) {
		t.Errorf("expected %s; got %s", v, b)
	}

	r.Flush()
}
