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
}

func New(secret, access, bucket, region string) S3 {
	return S3{secret: secret, access: access, bucket: bucket, region: region}
}

func (a S3) auth() aws.Auth {
	return aws.Auth{
		AccessKey: a.access,
		SecretKey: a.secret,
	}
}

func (a S3) remoteBucket() (*s3.Bucket, error) {
	if amazonConnection == nil {
		amazonConnection = s3.New(a.auth(), aws.Regions[a.region])
	}

	return amazonConnection.Bucket(a.bucket)
}

func (a S3) Put(uri string, data []byte) error {
	bucket, err := a.remoteBucket()

	if err != nil {
		return err
	}

	err = bucket.Put(uri, data, "application/json; charset=utf-8", s3.PublicRead)

	return err
}

func (a S3) Get(uri string) []byte {
	b := []byte{}

	bucket, err := a.remoteBucket()

	if err != nil {
		return []byte{}
	}

	b, err = bucket.Get(uri)

	if err != nil {
		return []byte{}
	}

	return b
}
