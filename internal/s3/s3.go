package s3

import (
	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
)

var amazonConnection *s3.S3

type S3 struct {
	access string
	secret string
	bucket string
	region string
	content string
}

func New(secret, access, bucket, region, content string) S3 {
	return S3{secret: secret, access: access, bucket: bucket, region: region, content: content}
}

func (store S3) auth() aws.Auth {
	return aws.Auth{
		AccessKey: store.access,
		SecretKey: store.secret,
	}
}

func (store S3) remoteBucket() (*s3.Bucket, error) {
	if amazonConnection == nil {
		amazonConnection = s3.New(store.auth(), aws.Regions[store.region])
	}

	return amazonConnection.Bucket(store.bucket)
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

func (store S3) Flush() {
	// Todo: Should only flush file created by this package. Maybe use a special index file?
}
