package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/kyani-inc/storage.v1"
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
		prefix := "test/storage"

		store, err = storage.S3(access, secret, bucket, region, content, prefix)
		// all keys will be surrounded by the prefix value and content
		// extension (if recognized) like: "test/storage/name.json"

	case "folder":
		store, err = storage.Folder(os.Getenv("FOLDER_STORAGE_PATH"))

	case "redis":
		store = storage.Redis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	case "memcached":
		hosts := strings.Split(os.Getenv("MEMCACHED_HOSTS"), ",")
		store = storage.Memcached(hosts)

	case "bolt":
		store, err = storage.Bolt(os.Getenv("BOLTDB_FILE_PATH"))

	case "dynamodb":
		secret := os.Getenv("AWS_SECRET_KEY")
		access := os.Getenv("AWS_ACCESS_KEY")
		region := os.Getenv("AWS_REGION")
		table := os.Getenv("DYNAMODB_TABLE")

		store, err = storage.DynamoDB(access, secret, region, table)

	case "elasticsearch":
		host := os.Getenv("ES_HOST")
		scheme := os.Getenv("ES_SCHEME")
		index := os.Getenv("ES_INDEX")
		namespace := os.Getenv("APP_NAME")
		awsSecret := os.Getenv("AWS_SECRET_KEY")
		awsKey := os.Getenv("AWS_ACCESS_KEY")

		store, err = storage.ElasticSearch(host, scheme, index, namespace, awsKey, awsSecret)

	default:
		store = storage.Local()
	}

	if err != nil {
		// Handle the error appropriately.
		// You may not want to panic in your application.
		panic(err.Error())
	}
}

func main() {
	// This is designed for XML or JSON documents but any []byte can be used.
	err := store.Put("name", []byte("John Doe"))

	if err != nil {
		fmt.Println("Error saving name", err.Error())
		return
	}

	data := store.Get("name") // []byte("John Doe")

	if data == nil {
		fmt.Println("Missing key: name")
		return
	}

	fmt.Printf("Hello, %s.\n", data) // Hello, John Doe.

	store.Delete("name") // remove "name"

	store.Flush() // get rid of everything
}
