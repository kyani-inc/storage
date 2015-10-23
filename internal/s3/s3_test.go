package s3_test

import (
	"os"
	"testing"

	"github.com/kyani-inc/storage/internal/s3"
)

var (
	access  = os.Getenv("AWS_ACCESS")
	secret  = os.Getenv("AWS_SECRET")
	bucket  = os.Getenv("AWS_BUCKET")
	region  = os.Getenv("AWS_REGION")
	content = "application/json; charset=utf-8"
)

func TestS3(t *testing.T) {
	t.Skip("Need a way for CI to test this.")

	k, v := "greetings.json", "Hello World"

	s, err := s3.New(access, secret, bucket, region, content)

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
}
