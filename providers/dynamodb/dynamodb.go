// DynamoDB storage abstraction layer
package dynamodb

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDB struct {
	Table string
	Conn  *dynamodb.DynamoDB
}

// New creates an instance of the DynamoDB struct for us.
func New(access string, secret string, region string, tableName string) (DynamoDB, error) {
	creds := credentials.NewStaticCredentials(access, secret, "")

	cfg := aws.Config{
		Region: aws.String(region),
	}

	cfg.WithCredentials(creds)

	ddb := DynamoDB{
		Table: tableName,
		Conn:  dynamodb.New(session.New(), &cfg),
	}

	return ddb, nil
}

// Put overwrites or creates as needed a new file based on the key.
func (ddb DynamoDB) Put(key string, data []byte) error {

	if len(data) == 0 {
		return errors.New("Can't put empty values")
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(ddb.Table),
		Item: map[string]*dynamodb.AttributeValue{
			"key": &dynamodb.AttributeValue{
				S: aws.String(key),
			},
			"value": &dynamodb.AttributeValue{
				B: data,
			},
		},
	}
	_, err := ddb.Conn.PutItem(params)
	return err
}

// Get attempts to grab data from dynamodb
func (ddb DynamoDB) Get(key string) []byte {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(ddb.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	}

	results, err := ddb.Conn.GetItem(params)

	if err != nil {
		return []byte{}
	}

	if val, ok := results.Item["value"]; ok {
		return val.B
	}
	return []byte{}
}

// Delete removes a file by key.
func (ddb DynamoDB) Delete(key string) {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(ddb.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	}
	ddb.Conn.DeleteItem(params)
}

// Flush removes the entire db table. This should be called with caution!
func (ddb DynamoDB) Flush() {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(ddb.Table), // Required
	}

	ddb.Conn.DeleteTable(params)
}
