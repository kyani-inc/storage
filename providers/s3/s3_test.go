package s3_test

import (
	"os"
	"testing"

	"github.com/kyani-inc/storage/providers/s3"
	"github.com/subosito/gotenv"
)

var (
	access, secret, bucket, region string

	prefix  = "test/storage"
	content = "application/json; charset=utf-8"
)

func init() {
	gotenv.Load(".env")

	access = os.Getenv("AWS_ACCESS")
	secret = os.Getenv("AWS_SECRET")
	bucket = os.Getenv("AWS_BUCKET")
	region = os.Getenv("AWS_REGION")
}

func emptyVars() bool {
	return access == "" || secret == "" || bucket == "" || region == ""
}

func TestS3(t *testing.T) {
	if emptyVars() {
		t.Skip("need AWS credentials in order to test")
	}

	k, v := "greetings.json", "Hello World"

	s, err := s3.New(access, secret, bucket, region, content, prefix)

	if err != nil {
		t.Fatal(err.Error())
	}

	err = s.Put(k, []byte(v))

	if err != nil {
		t.Fatal(err.Error())
	}

	b := s.Get(k)

	if v != string(b) {
		t.Errorf("expected %s; got %s", v, b)
	}

	s.Delete(k)

	c := s.Get(k)

	if len(c) > 1 {
		t.Errorf("expected empty return; got %s", c)
	}

	s.Flush()
}
