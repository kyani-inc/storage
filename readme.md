# Storage [![Build Status](https://travis-ci.org/kyani-inc/storage.svg)](https://travis-ci.org/kyani-inc/storage)&nbsp;[![godoc reference](https://godoc.org/github.com/kyani-inc/storage?status.png)](https://godoc.org/github.com/kyani-inc/storage)

Unified key/value storage interface for several backing technologies.

- **Bolt** ([BoltDB](https://github.com/boltdb/bolt) backed store; production ready)
- **Memcached** ([Memcached](http://memcached.org/) backed store; production ready)
- **Redis** ([Redis](http://redis.io/) backed store; production ready)
- **S3** ([AWS S3](https://aws.amazon.com/s3/) backed store; production ready)
- **Folder** (local folder backed store; intended for dev)
- **Local** (application's memory backed store; intended for dev)

# Usage and Examples

`go get https://godoc.org/gopkg.in/kyani-inc/storage.v1`

You can write your application to use the storage interface and then change out the backing technology based on 
environment and application needs. For example, running `Redis` in production and using `Local` or `Folder` for 
local development.

```go
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

	fmt.Printf("Hello, %s.\n", data) // Hello, John Doe.

	store.Delete("name") // remove "name"

	store.Flush() // get rid of everything
}
```

# Testing

run `go test -v ./...`

## Env Vars for Testing

To test certain drivers you will need to provide the services to test against. The tests are built using 
[gotenv](https://github.com/subosito/gotenv) so that you can provide environment variables for the tests.
Otherwise tests that require these variables will be skipped.

For example: Testing the Redis implementation on your local machine would require the file `internal/redis/.env` with the following
contents (replace HOST and PORT with your values).

```
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
```

