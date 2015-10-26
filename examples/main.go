package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kyani-inc/storage"
)

var store storage.Storage

func init() {
	var err error

	switch os.Getenv("STORAGE_ENGINE") {
	case "s3":
		secret := os.Getenv("AWS_SECRET_KEY")
		access := os.Getenv("AWS_ACCESS_KEY")
		bucket := os.Getenv("AWS_BUCKET")
		region := os.Getenv("AWS_REGION")
		content := "application/json; charset=utf-8"

		store, err = storage.S3(secret, access, bucket, region, content)

	case "folder":
		store, err = storage.Folder(os.Getenv("FOLDER_STORAGE_PATH"))

	case "redis":
		store = storage.Redis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	case "memcache":
		hosts := strings.Split(os.Getenv("MEMCACHE_HOSTS"), ",")
		store = storage.Memcache(hosts)

	case "bolt":
		store, err = storage.Bolt(os.Getenv("BOLTDB_FILE_PATH"))

	default:
		store = storage.Local()
	}

	if err != nil {
		panic(err.Error())
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
