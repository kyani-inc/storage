// Package storage creates an interface for several key/value storage technologies.
package storage

import (
	"github.com/kyani-inc/storage/providers/bolt"
	"github.com/kyani-inc/storage/providers/dynamodb"
	"github.com/kyani-inc/storage/providers/folder"
	"github.com/kyani-inc/storage/providers/local"
	"github.com/kyani-inc/storage/providers/memcached"
	"github.com/kyani-inc/storage/providers/redis"
	"github.com/kyani-inc/storage/providers/s3"
)

// Storage represents a storage facility agnostic of the backing technology.
type Storage interface {
	// Get returns data by key. Missing keys return empty []byte.
	Get(key string) []byte

	// Put will overwrite data by key.
	Put(key string, data []byte) error

	// Delete will attempt to remove a value by key. Idempotent.
	Delete(key string)

	// Flush clears all keys from the storage. Idempotent.
	Flush()
}

// Bolt utilizes a BoltDB database (https://github.com/boltdb/bolt) for storage.
func Bolt(path string) (Storage, error) {
	return bolt.New(path)
}

// Local uses the applications memory for storage.
func Local() Storage {
	return local.New()
}

// Folder uses the application's underlying file structure for storage.
func Folder(path string) (Storage, error) {
	return folder.New(path)
}

// S3 uses Amazon AWS S3 for storage.
//
// Every key combined with the prefix (to support Flush) and a possible extension determined by content.
// The field content is the content type to use with all keys. For example: "application/json; charset=utf-8".
func S3(access, secret, bucket, region, content, prefix string) (Storage, error) {
	return s3.New(access, secret, bucket, region, content, prefix)
}

// Redis uses a Redis instance for storage.
func Redis(host, port string) Storage {
	return redis.New(host, port)
}

// Memcache uses one or more Memcache hosts for storage.
func Memcached(hosts []string) Storage {
	return memcached.New(hosts)
}

// AWS DynamoDB storage.
func DynamoDB(region string, endpoint string, table string) (Storage, error) {
	return dynamodb.New(region, endpoint, table)
}
