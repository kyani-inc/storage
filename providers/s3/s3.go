// s3 handles storing documents in AWS S3 buckets.
package s3

import (
	"errors"
	"fmt"
	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
)

// S3 represents the data needed to communicate with AWS S3 service.
type S3 struct {
	access    string
	secret    string
	bucket    string
	region    string
	content   string
	prefix    string
	extension string
	conn      *s3.S3
}

// New creates a connection to an AWS S3 service bucket.
//
// The parameters access, secret, bucket, and region are issued from AWS.
//
// The parameter content should be a recognized Content-Type to help download applications correctly
// associate the file contents. AWS allows this to be blank.
//
// The parameter prefix is required so that Flush is a simpler process. Deleting and recreating the
// bucket is not supported.
func New(access, secret, bucket, region, content, prefix string) (S3, error) {
	s := S3{
		access:  access,
		secret:  secret,
		bucket:  bucket,
		region:  region,
		content: content,
		prefix:  prefix,
	}

	if err := s.validate(); err != nil {
		return s, err
	}

	s.extension = parseExt(s.content)

	_, err := s.remoteBucket()

	return s, err
}

// Put overwrites or creates as needed a new file based on the key.
func (store S3) Put(key string, data []byte) error {
	bucket, err := store.remoteBucket()

	if err != nil {
		return err
	}

	err = bucket.Put(store.uri(key), data, store.content, s3.PublicRead)

	return err
}

// Get attempts to access the contents of a file by key.
func (store S3) Get(key string) []byte {
	var b []byte

	bucket, err := store.remoteBucket()

	if err != nil {
		return b
	}

	b, err = bucket.Get(store.uri(key))

	if err != nil {
		return nil
	}

	return b
}

// Delete removes a file by key.
func (store S3) Delete(key string) {
	bucket, err := store.remoteBucket()

	if err != nil {
		return
	}

	bucket.Del(store.uri(key))
}

// Flush removes the prefix folder from S3, effectively removing all files/keys.
func (store S3) Flush() {
	bucket, err := store.remoteBucket()

	if err != nil {
		return
	}

	if store.prefix != "" {
		bucket.Del(store.prefix)
	}
}

func (store S3) auth() aws.Auth {
	return aws.Auth{
		AccessKey: store.access,
		SecretKey: store.secret,
	}
}

func (store S3) uri(v string) string {
	return store.prefix + v + store.extension
}

func (store S3) remoteBucket() (*s3.Bucket, error) {
	if store.conn == nil {
		store.conn = s3.New(store.auth(), aws.Regions[store.region])
	}

	return store.conn.Bucket(store.bucket)
}

func (store S3) validate() error {
	errs := []string{}

	if store.access == "" {
		errs = append(errs, "access key missing")
	}

	if store.secret == "" {
		errs = append(errs, "secret key missing")
	}

	if store.region == "" {
		errs = append(errs, "region missing")
	}

	if store.bucket == "" {
		errs = append(errs, "bucket missing")
	}

	if store.prefix == "" {
		errs = append(errs, "prefix missing")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("S3 errors %+v", errs))
}
