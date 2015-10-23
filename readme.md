# Storage [![Build Status](https://travis-ci.org/kyani-inc/storage.svg)](https://travis-ci.org/kyani-inc/storage)&nbsp;[![godoc reference](https://godoc.org/github.com/kyani-inc/storage?status.png)](https://godoc.org/github.com/kyani-inc/storage)

Unified key/value storage interface for several backing technologies.

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
```

# ToDo

- [x] Folder Support
- [x] Local Memory Support
- [x] Memcache Support
- [x] Redis Support
- [x] S3 Support
- [x] Basic Example
- [ ] S3 support of `storage.Flush()` [see here](s3/s3.go#L68)
