// DynamoDB storage abstraction layer
package dynamodb

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDB struct {
	Table string
	Conn  *dynamodb.DynamoDB
}

// New creates an instance of the DynamoDB struct for us.
func New(region string, endPoint string, tableName string) (DynamoDB, error) {
	cfg := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endPoint),
	}

	ddb := DynamoDB{
		Table: tableName,
		Conn:  dynamodb.New(session.New(), &cfg),
	}

	err := ddb.createTable()

	if err != nil {
		return ddb, errors.New("DB Error")
	}

	return ddb, nil
}

// tableExists checks to see if it can query the table meta data.
// if not, the table *probably* doesn't exists.
func (ddb DynamoDB) tableExists() bool {
	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(ddb.Table), // Required
	}

	resp, err := ddb.Conn.DescribeTable(params)

	if err != nil {
		// We'll assume just this is a "ResourceNotFoundException".
		return false
	}

	if resp != nil && resp.Table != nil {
		status := resp.Table.TableStatus

		// yeah...you have to dereference a string -_-
		if *status == "ACTIVE" {
			return true
		}
	}

	return false
}

// createTable yes, this is all required.
func (ddb DynamoDB) createTable() (err error) {
	exists := ddb.tableExists()

	if exists == true {
		return nil
	}

	table := &dynamodb.CreateTableInput{
		TableName: aws.String(ddb.Table),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(2),
		},
	}

	resp, err := ddb.Conn.CreateTable(table)

	if err != nil {
		// Failed to create table
		return err
	}

	if resp != nil && resp.TableDescription != nil {
		status := resp.TableDescription.TableStatus

		if *status == "ACTIVE" {
			// The table is ready for use..
		}

		if *status == "CREATING" {
			// Block for a bit to ensure the table exists before we attempt to start using it...
			// we may want to make sure we create the table in the AWS UI before using it.
			time.Sleep(15 * time.Second)
		}
	}

	return nil
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
	_, _ = ddb.Conn.DeleteItem(params)
}

// Flush removes the entire db table. This should be called with caution!
func (ddb DynamoDB) Flush() {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(ddb.Table), // Required
	}

	_, err := ddb.Conn.DeleteTable(params)

	if err != nil {
		fmt.Println("Failed to flush table", err.Error())
	}
}
