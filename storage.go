// Package storage creates an interface for several storage technologies.
package storage

import (
	"github.com/kyani-inc/storage/internal/bolt"
	"github.com/kyani-inc/storage/internal/folder"
	"github.com/kyani-inc/storage/internal/local"
	"github.com/kyani-inc/storage/internal/memcache"
	"github.com/kyani-inc/storage/internal/redis"
	"github.com/kyani-inc/storage/internal/s3"
)

// Storage represents a storage facility agnostic of the backing technology.
type Storage interface {
	// Get returns data by key. Missing keys return empty []byte.
	Get(key string) []byte

	// Put will overwrite data by key.
	Put(key string, data []byte) error

	// Flush clears all keys from the storage.
	Flush()
}

// Bolt utilizes a boltdb database for storage.
func Bolt(path string) (Storage, error) {
	return bolt.New(path), nil
}

// Local uses the applications memory for storage.
func Local() Storage {
	return local.New()
}

// Folder uses the application's underlying file structure for storage.
// Note: Flush() will delete path.
func Folder(path string) (Storage, error) {
	return folder.New(path)
}

// S3 uses Amazon AWS S3 for storage.
// Every key is treated as a URI and makes a new file on first Put.
// The field content is the content type to use with all keys. For example: "application/json; charset=utf-8".
func S3(secret, access, bucket, region, content string) (Storage, error) {
	return s3.New(secret, access, bucket, region, content)
}

// Redis uses a Redis instance for storage.
func Redis(host, port string) (Storage, error) {
	return redis.New(host, port), nil
}

// Memcache uses one or more Memcache hosts for storage.
func Memcache(hosts []string) (Storage, error) {
	return memcache.New(hosts), nil
}

// Basic is deprecated and is here for backwards compatibility. Use Local().
func Basic() Storage {
	return Local()
}
