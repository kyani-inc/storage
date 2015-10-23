package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kyani-inc/storage"
)

var store storage.Storage

func init() {
	switch os.Getenv("StorageEngine") {
	case "s3":
		secret := os.Getenv("AWS_SECRET_KEY")
		access := os.Getenv("AWS_ACCESS_KEY")
		bucket := os.Getenv("AWS_BUCKET")
		region := os.Getenv("AWS_REGION")
		content := "application/json; charset=utf-8"

		store = storage.S3(secret, access, bucket, region, content)

	case "folder":
		store = storage.Folder("/tmp/storage")

	case "redis":
		store = storage.Redis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	case "memcache":
		hosts := strings.Split(os.Getenv("MEMCACHE_HOSTS"), ",")
		store = storage.Memcache(hosts)

	case "bolt":
		store = storage.Bolt("/tmp/storage.db")

	default:
		store = storage.Local()
	}
}

func main() {
	err := store.Put("mydata", []byte("Hello World"))

	if err != nil {
		fmt.Println("Error saving mydata", err.Error())
		return
	}

	data := store.Get("mydata")

	fmt.Println("Got data", string(data))
}
