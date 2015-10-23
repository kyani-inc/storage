package storage

import (
	"github.com/kyani-inc/storage/folder"
	"github.com/kyani-inc/storage/local"
	"github.com/kyani-inc/storage/redis"
	"github.com/kyani-inc/storage/s3"
)

type Storage interface {
	Get(key string) []byte
	Put(key string, data []byte) error
}

func Basic() Storage {
	return local.New()
}

func Folder(path string) Storage {
	if path == "" {
		return local.New()
	}

	return folder.New(path)
}

func S3(secret, access, bucket, region string) Storage {
	return s3.New(secret, access, bucket, region)
}

func Redis(host, port string) Storage {
	if host == "" || port == "" {
		return local.New()
	}

	return redis.New(host, port)
}
