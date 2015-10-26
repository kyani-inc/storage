# Storage [![Build Status](https://travis-ci.org/kyani-inc/storage.svg)](https://travis-ci.org/kyani-inc/storage)&nbsp;[![godoc reference](https://godoc.org/github.com/kyani-inc/storage?status.png)](https://godoc.org/github.com/kyani-inc/storage)

Unified key/value storage interface for several backing technologies.

- **Bolt** ([BoltDB](https://github.com/boltdb/bolt) file backed; production ready)
- **Folder** (can be slow due to locking; useful for dev or testing)
- **Local** (application memory; useful for dev or testing)
- **Memcache** (production ready)
- **Redis** (production ready)
- **S3** (production ready)

# Usage and Examples

You can write your application to use the storage interface and then change out the backing technology based on 
environment and application needs. For example, running `Redis` in production and using `Local` or `Folder` for 
local development.

```go
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
}

```

# To Do

- [x] Folder Support
- [x] Local Memory Support
- [x] Memcache Support
- [x] Redis Support
- [x] S3 Support
- [x] Basic Example
- [ ] S3 support of `storage.Flush()` [see here](s3/s3.go#L68)
