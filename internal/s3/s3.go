package s3

import (
	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
)

type S3 struct {
	access string
	secret string
	bucket string
	region string
	content string
	conn *s3.S3
}

func New(access, secret, bucket, region, content string) (S3, error){
	s := S3{secret: secret, access: access, bucket: bucket, region: region, content: content}

	_, err := s.remoteBucket()

	return s, err
}

func (store S3) auth() aws.Auth {
	return aws.Auth{
		AccessKey: store.access,
		SecretKey: store.secret,
	}
}

func (store S3) remoteBucket() (*s3.Bucket, error) {
	if store.conn == nil {
		store.conn = s3.New(store.auth(), aws.Regions[store.region])
	}

	return store.conn.Bucket(store.bucket)
}

func (store S3) Put(uri string, data []byte) error {
	bucket, err := store.remoteBucket()

	if err != nil {
		return err
	}

	err = bucket.Put(uri, data, store.content, s3.PublicRead)

	return err
}

func (store S3) Get(uri string) []byte {
	b := []byte{}

	bucket, err := store.remoteBucket()

	if err != nil {
		return []byte{}
	}

	b, err = bucket.Get(uri)

	if err != nil {
		return []byte{}
	}

	return b
}

func (store S3) Delete(uri string) {
	bucket, err := store.remoteBucket()

	if err != nil {
		return
	}

	bucket.Del(uri)
}

func (store S3) Flush() {
	// ToDo: Should only flush file created by this package. Maybe use a special index file?
}
